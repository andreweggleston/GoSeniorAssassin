package socket

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/andreweggleston/GoSeniorAssassin/routes/socket"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers/hooks"
)

func RegisterHandlers() {
	socket.AuthServer.OnDisconnect = hooks.OnDisconnect
	socket.UnauthServer.OnDisconnect = func(string, *jwt.Token) {pprof.Clients.Add(-1)}
}