package admin

import (
	"html/template"
	"net/http"
	"golang.org/x/net/xsrftoken"
	"github.com/Sirupsen/logrus"
	"github.com/andreweggleston/GoSeniorAssassin/config"
)

var banForm = map[string]string{
	"chat":            "Chatting",
	"full":            "Full ban",
}

var roleForm = map[string]string{
	"admin": "Add Administrator",
	"mod":   "Add Moderator",
}

var adminPageTempl *template.Template

func ServeAdminPage(w http.ResponseWriter, r *http.Request) {
	err := adminPageTempl.Execute(w, map[string]interface{}{
		"BanForms":  banForm,
		"RoleForms": roleForm,
		"XSRFToken": xsrftoken.Generate(config.Constants.CookieStoreSecret, "admin", "POST"),
	})
	if err != nil {
		logrus.Error(err)
	}
}