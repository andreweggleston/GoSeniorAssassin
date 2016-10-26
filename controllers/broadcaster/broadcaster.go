package broadcaster

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/sessions"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/andreweggleston/GoSeniorAssassin/routes/socket"
)

func SendMessage(id string, event string, content interface{}) {
	sockets, ok := sessions.GetSockets(id)
	if !ok {
		return
	}

	for _, socket := range sockets {
		go func(so *wsevent.Client) {
			so.EmitJSON(helpers.NewRequest(event, content))
		}(socket)
	}

	return
}

func SendMessageToRoom(r string, event string, content interface{}) {
	v := helpers.NewRequest(event, content)

	socket.AuthServer.BroadcastJSON(r, v)
	socket.UnauthServer.BroadcastJSON(r, v)
}

func SendMessageSkipIDs(skipID, id, event string, content interface{}) {
	sockets, ok := sessions.GetSockets(id)
	if !ok {
		return
	}

	for _, socket := range sockets {
		if socket.ID != skipID {
			socket.EmitJSON(helpers.NewRequest(event, content))
		}
	}
}
