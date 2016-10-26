package socket

import (
	"github.com/TF2Stadium/wsevent"
	"errors"
	"github.com/andreweggleston/GoSeniorAssassin/routes/socket/middleware"
)

var (
	//AuthServer is the wsevent server where authenticated (logged in) users/sockets
	//are added to
	AuthServer = wsevent.NewServer(middleware.JSONCodec{}, func(_ *wsevent.Client, _ struct{}) interface{} {
		return errors.New("No such request.")
	})
	//UnauthServer is the wsevent server where unauthenticated users/sockets
	//are added to
	UnauthServer = wsevent.NewServer(middleware.JSONCodec{}, func(_ *wsevent.Client, _ struct{}) interface{} {
		return errors.New("You aren't logged in.")
	})
)

// Wait for all websocket requests to complete
func Wait() {
	AuthServer.Requests.Wait()
	UnauthServer.Requests.Wait()
}