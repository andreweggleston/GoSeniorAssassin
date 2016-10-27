package hooks

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/sessions"
	"time"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

func AfterConnect(server *wsevent.Server, so *wsevent.Client) {
	server.Join(so, "0_public")
}

var emptyMap = make(map[string]string)

func AfterConnectLoggedIn(so *wsevent.Client, player *player.Player) {
	sessions.AddSocket(player.ID, so)

	if time.Since(player.ProfileUpdatedAt) >= 30*time.Minute {
		player.UpdatePlayerInfo
	}

	if player.Settings != nil {
		so.EmitJSON(helpers.NewRequest("playerSettings", player.Settings))
	} else {
		so.EmitJSON(helpers.NewRequest("playerSettings", emptyMap))
	}
}