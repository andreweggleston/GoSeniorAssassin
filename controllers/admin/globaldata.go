package admin

import (
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/databaseAssassin"
	"fmt"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
)

func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	values := r.Form

	announcement := values.Get("announcement")

	databaseAssassin.DB.Exec("UPDATE players SET announcement='" + announcement + "';")

	controllerhelpers.SendGlobalMessage("[Senior Assassin] New Announcement" + announcement)

	fmt.Fprintf(w, "Global Announcement updated to %s", announcement)
}

func UpdateSafetyItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	values := r.Form

	safetyItem := values.Get("safetyitem")

	databaseAssassin.DB.Exec("UPDATE players SET safety_item='" + safetyItem + "';")

	controllerhelpers.SendGlobalMessage("[Senior Assassin] Safety Item update: " + safetyItem)

	fmt.Fprintf(w, "Global Safety Item updated to %s", safetyItem)
}

func UpdateKillByDate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	values := r.Form

	killDate := values.Get("killdate")

	databaseAssassin.DB.Exec("UPDATE players SET kill_by_date='" + killDate + "';")

	controllerhelpers.SendGlobalMessage("[Senior Assassin] Round End date updated: " + killDate)

	fmt.Fprintf(w, "Global Kill By Date updated to %s", killDate)
}