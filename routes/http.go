package routes

import (
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/controllers"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/login"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/admin"
	"github.com/andreweggleston/GoSeniorAssassin/config"
)


type route struct {
	pattern string
	handler http.HandlerFunc
}

var httpRoutes = []route{
	{"/", controllers.MainHandler},
	{"/websocket/", controllers.SocketHandler},

	{"/login", login.GoogleLoginHandler},
	{"/auth", login.GoogleAuthHandler},
	{"/logout", login.GoogleLogoutHandler},


	{"/admin", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.ServeAdminPage)},
	{"/admin/roles", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.ChangeRole)},
	{"/admin/updateannouncement", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.UpdateAnnouncement)},
	{"/admin/updatesafety", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.UpdateSafetyItem)},
	{"/admin/updatekilldate", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.UpdateKillByDate)},
	{"/admin/chatlogs", chelpers.FilterHTTPRequest(helpers.ActionViewLogs, admin.GetChatLogs)},
	{"/admin/banlogs", chelpers.FilterHTTPRequest(helpers.ActionViewLogs, admin.GetBanLogs)},
}

func SetupHTTP(mux *http.ServeMux) {
	for _, httpRoute := range httpRoutes {
		mux.HandleFunc(httpRoute.pattern, httpRoute.handler)
	}
	mux.Handle("/demos/", http.StripPrefix("/demos/", http.FileServer(http.Dir("/"))))
	//mux.Handle("/oauth2callback", http.HandlerFunc(redirectToHTTP))
	//mux.Handle("/startBNetLogin", http.HandlerFunc(redirectToHTTP))

	if config.Constants.ServeStatic {
		mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "views/static.html")
		})

	}
}