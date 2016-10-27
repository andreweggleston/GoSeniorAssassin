package admin

import (
	"net/http"
	"golang.org/x/net/xsrftoken"
	"fmt"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)

func ChangeRole(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	values := r.Form
	studentid := values.Get("studentid")
	remove := values.Get("remove")
	token := values.Get("xsrf-token")
	if !xsrftoken.Valid(token, config.Constants.CookieStoreSecret, "admin", "POST") {
		http.Error(w, "invalid xsrf token", http.StatusBadRequest)
		return
	}

	role, ok := map[string]authority.AuthRole{
		"admin": helpers.RoleAdmin,
		"mod":   helpers.RoleMod,
	}[values.Get("role")]
	if !ok {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	player, err := player.GetPlayerByStudentID(studentid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if remove == "true" {
		player.Role = 0
		player.Save()
		fmt.Fprintf(w, "Player %s (%s) has been removed as %s", player.Name, player.ID, helpers.RoleNames[role])
		return
	}

	player.Role = role
	player.Save()
	fmt.Fprintf(w, "Player %s (%s) has been made a %s", player.Name, player.ID, helpers.RoleNames[role])
	return
}

func Remove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	studentid := r.Form.Get("studentid")
	player, err := player.GetPlayerByStudentID(studentid)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	player.Role = authority.AuthRole(0)
	player.Save()
	fmt.Fprintf(w, "%s (%s) is no longer an admin/mod", player.Name, player.ID)
}