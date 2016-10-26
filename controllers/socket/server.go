package socket

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/andreweggleston/GoSeniorAssassin/routes/socket"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers/hooks"
	"github.com/andreweggleston/GoSeniorAssassin/internal/pprof"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/handler"
)

func RegisterHandlers() {
	socket.AuthServer.OnDisconnect = hooks.OnDisconnect
	socket.UnauthServer.OnDisconnect = func(string, *jwt.Token) {pprof.Clients.Add(-1)}

	socket.AuthServer.Register(handler.Global{})
	socket.AuthServer.Register(handler.Player{})
	socket.AuthServer.Register(handler.Chat{})
	
	socket.UnauthServer.Register(handler.Unauth{})
}