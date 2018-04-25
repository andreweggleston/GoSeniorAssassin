package hooks

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket/sessions"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/sirupsen/logrus"
)

func AfterConnect(server *wsevent.Server, so *wsevent.Client) {
	server.Join(so, "0_public")
}

var emptyMap = make(map[string]string)

func AfterConnectLoggedIn(so *wsevent.Client, player *player.Player) {
	sessions.AddSocket(player.StudentID, so)

	err := player.UpdatePlayerData()
	if err != nil{
		logrus.Error(err)
	}

	if player.Settings != nil {
		so.EmitJSON(helpers.NewRequest("playerSettings", player.Settings))
	} else {
		so.EmitJSON(helpers.NewRequest("playerSettings", emptyMap))
	}

	player.SetPlayerProfile()
	so.EmitJSON(helpers.NewRequest("playerProfile", player))
}