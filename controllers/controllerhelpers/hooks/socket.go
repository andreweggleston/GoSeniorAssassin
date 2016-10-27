package hooks

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/sessions"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
)

func OnDisconnect(socketID string, token *jwt.Token) {
	if token != nil {
		//player was logged in
		player := chelpers.GetPlayer(token)
		if player == nil {
			return
		}

		//TODO: get the lobby shit out

		sessions.RemoveSocket(socketID, player.StudentID)

	}
}

