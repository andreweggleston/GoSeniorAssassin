package controllers

import (
	"errors"
	"github.com/TF2Stadium/wsevent"
	"fmt"
	"net/http"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/routes/socket"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers/hooks"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	token, err := chelpers.GetToken(r)
	if err != nil && err != http.ErrNoCookie { //invalid jwt token
		token = nil
	}

	//check if player is in the whitelist
	if config.Constants.IDWhitelist != "" {
		if token == nil {
			// player isn't logged in,
			// and access is restricted to logged in people
			http.Error(w, "Not logged in", http.StatusForbidden)
			return
		}

		if !chelpers.IsIDWhitelisted(token.Claims["student_id"].(string)) {
			http.Error(w, "you're not in the beta", http.StatusForbidden)
			return
		}
	}

	var so *wsevent.Client

	if token != nil { //received valid jwt
		so, err = socket.AuthServer.NewClient(upgrader, w, r)
	} else {
		so, err = socket.UnauthServer.NewClient(upgrader, w, r)
	}

	if err != nil {
		return
	}

	so.Token = token

	//logrus.Debug("Connected to Socket")
	err = SocketInit(so)
	if err != nil {
		logrus.Error(err)
		so.Close()
		return
	}
}

var ErrRecordNotFound = errors.New("Player record for found.")

//SocketInit initializes the websocket connection for the provided socket
func SocketInit(so *wsevent.Client) error {
	loggedIn := so.Token != nil

	if loggedIn {
		hooks.AfterConnect(socket.AuthServer, so)

		studentid := so.Token.Claims["student_id"].(string)

		player, err := player.GetPlayerByStudentID(studentid)
		if err != nil {
			return fmt.Errorf("Couldn't find player record for %s", studentid)
		}

		hooks.AfterConnectLoggedIn(so, player)
	} else {
		hooks.AfterConnect(socket.UnauthServer, so)
		so.EmitJSON(helpers.NewRequest("playerSettings", "{}"))
		so.EmitJSON(helpers.NewRequest("playerProfile", "{}"))
	}

	so.EmitJSON(helpers.NewRequest("socketInitialized", "{}"))

	return nil
}