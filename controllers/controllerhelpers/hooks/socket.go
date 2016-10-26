package hooks

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/sessions"
)

func OnDisconnect(socketID string, token *jwt.Token) {
	if token != nil { //player was logged in
		player := chelpers.GetPlayer(token)
		if player == nil {
			return
		}

		//TODO: get the lobby shit out

		sessions.RemoveSocket(socketID, player.ID)
		if id != 0 {
			lob, _ := lobby.GetLobbyByID(id)
			lob.RemoveSpectator(player, true)
		}

		id, _ = player.GetLobbyID(true)
		//if player is in a waiting lobby, and hasn't connected for > 30 seconds,
		//remove him from it. Here, connected = player isn't connected from any tab/window
		if id != 0 && sessions.ConnectedSockets(player.SteamID) == 0 {
			sessions.AfterDisconnectedFunc(player.SteamID, time.Second*30, func() {
				lob, _ := lobby.GetLobbyByID(id)
				if lob.State == lobby.Waiting {
					lob.RemovePlayer(player)
				}
			})
		}
	}

