package handler

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)

type Unauth struct {}

func (Unauth) Name(s string) string {
	return string((s[0])+32) + s[1:]
}

func (Unauth) PlayerProfile(so *wsevent.Client, args struct {
	Steamid *string `json:"steamid"`
}) interface{} {

	player, err := player.GetPlayerByID(*args.Steamid)
	if err != nil {
		return err
	}

	player.SetPlayerProfile()
	return newResponse(player)
}
