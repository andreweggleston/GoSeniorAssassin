package controllerhelpers

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
)

func BroadcastScrollback(so *wsevent.Client, room uint) {
	messages, err := chat.GetScrollback(int(room))
	if err != nil {
		return
	}

	so.EmitJSON(helpers.NewRequest("chatScrollback", messages))
}
