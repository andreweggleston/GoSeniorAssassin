package login

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"encoding/base64"
	"crypto/rand"

	"io/ioutil"
	"log"
	"fmt"
	"net/http"

	"github.com/andreweggleston/GoSeniorAssassin/config"
	"time"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"strings"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
)

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

var (
	conf = &oauth2.Config{
		ClientID:     "812567846854-k544u30oiihu7uaqo8b8gs8iup50od14.apps.googleusercontent.com",
		ClientSecret: "i7j027VP4RJLYBiE1jALdSo9",
		RedirectURL:  "http://10.0.0.5.nip.io:8081/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	oauthStateString = "randomstate"
	cred             Credentials
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(oauthStateString)
	logrus.Info(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the exchange code to initiate a transport.
	retrievedState := r.FormValue("state")
	if retrievedState != oauthStateString {
		fmt.Printf("oauth state is invalid, expected '%s', got '%s' \n", oauthStateString, retrievedState)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("r.ParseForm() failed with %s\n", err)
	}

	code := r.Form.Get("code")

	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("conf.Exchange() failed with '%s' \n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Printf("client.Get() failed with '%s' \n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	log.Println(string(data))

	user := ParseBody(data)

	studentid := strings.Split(user.Email, "@student.wayland.k12.ma.us")[0]

	if !strings.Contains(user.Email, "@student.wayland.k12.ma.us") {
		http.Error(w, "Sign up with your school account!", http.StatusForbidden)
		http.Redirect(w, r, "/", http.StatusForbidden)
	} else {
		p, err := player.GetPlayerByStudentID(studentid)
		if err != nil {

			p, err = player.NewPlayer(studentid)

			if err != nil {
				logrus.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			p.Name = strings.Replace(studentid, "_", " ", -1)
			log.Print(strings.Split(user.Email, "@student.wayland.k12.ma.us")[0])
			p.Sub = user.Sub
			p.Email = user.Email

			database.DB.Create(p)
		}


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

		http.Redirect(w, r, config.Constants.LoginRedirectPath, 303)
	}
}

func ParseBody(body []byte) (*User) {
	var s = new(User)
	err := json.Unmarshal(body, &s)
	if err != nil {
		logrus.Fatal(err)
	}
	return s
}

func GoogleLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth-jwt")
	if err != nil { //idiot wasnt even logged in LUL
		return
	}

	cookie.Domain = config.Constants.CookieDomain
	cookie.MaxAge = -1
	cookie.Expires = time.Time{}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, config.Constants.LoginRedirectPath, 303)
}
