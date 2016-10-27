package login

import (
	"time"
	"github.com/yohcop/openid-go"
	"net/url"
	"github.com/Sirupsen/logrus"
	"regexp"
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
)

var (
	nonceStore     = openid.NewSimpleNonceStore()
	discoveryCache = openid.NewSimpleDiscoveryCache()
)

func SteamLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirecturl, _ := url.Parse(config.Constants.PublicAddress)
	redirecturl.Path = "openidcallback"

	referer, ok := r.Header["Referer"]
	if ok {
		values := redirecturl.Query()
		values.Set("referer", referer[0])
		redirecturl.RawQuery = values.Encode()
	}

	//if url, err := openid.RedirectURL("http://steamcommunity.com/openid",
	//	redirecturl.String(),
	//	config.Constants.OpenIDRealm); err == nil {
	//	http.Redirect(w, r, url, 303)
	//} else {
	//	logrus.Error(err)
	//}
}

func SteamMockLoginHandler(w http.ResponseWriter, r *http.Request) {
	if !config.Constants.MockupAuth {
		http.NotFound(w, r)
		return
	}

	studentid := r.URL.Query().Get("studentid")
	if studentid == "" {
		http.Error(w, "No StudentID given", http.StatusBadRequest)
		return
	}

	p, err := player.GetPlayerByStudentID(studentid)
	if err != nil {
		p, err = player.NewPlayer(studentid)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		database.DB.Create(p)
	}

	p.UpdatePlayerInfo()
	key := controllerhelpers.NewToken(p)
	cookie := &http.Cookie{
		Name:    "auth-jwt",
		Value:   key,
		Path:    "/",
		Domain:  config.Constants.CookieDomain,
		Expires: time.Now().Add(30 * 24 * time.Hour),
		//Secure: true,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, config.Constants.LoginRedirectPath, 303)
}

func SteamLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth-jwt")
	if err != nil { //user wasn't even logged in ಠ_ಠ
		return
	}

	cookie.Domain = config.Constants.CookieDomain
	cookie.MaxAge = -1
	cookie.Expires = time.Time{}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, config.Constants.LoginRedirectPath, 303)
}

var reSteamID = regexp.MustCompile(`http://steamcommunity.com/openid/id/(\d+)`)

func SteamLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	refererURL := r.URL.Query().Get("referer")

	publicURL, _ := url.Parse(config.Constants.PublicAddress)
	// this wouldnt be used anymore, so modify it directly
	r.URL.Scheme = publicURL.Scheme
	r.URL.Host = publicURL.Host
	idURL, err := openid.Verify(r.URL.String(), discoveryCache, nonceStore)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parts := reSteamID.FindStringSubmatch(idURL)
	logrus.Info(parts)
	if len(parts) != 2 {
		http.Error(w, "Authentication failed, please try again.", 500)
		return
	}

	studentid := parts[1]

	if config.Constants.IDWhitelist != "" &&
		!controllerhelpers.IsIDWhitelisted(studentid) {
		//Use a more user-friendly message later
		http.Error(w, "Sorry, you're not in the closed alpha.", 403)
		return
	}

	p, err := player.GetPlayerByStudentID(studentid)
	if err != nil {
		p, err = player.NewPlayer(studentid)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		database.DB.Create(p)
	}

	go func() {
		if time.Since(p.ProfileUpdatedAt) >= 1*time.Hour {
			p.UpdatePlayerInfo()
		}
	}()

	key := controllerhelpers.NewToken(p)
	cookie := &http.Cookie{
		Name:     "auth-jwt",
		Value:    key,
		Path:     "/",
		Domain:   config.Constants.CookieDomain,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   config.Constants.SecureCookies,
	}

	http.SetCookie(w, cookie)
	if refererURL != "" {
		http.Redirect(w, r, refererURL, 303)
		return
	}

	http.Redirect(w, r, config.Constants.LoginRedirectPath, 303)
}
