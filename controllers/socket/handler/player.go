package handler

import (
	"github.com/TF2Stadium/wsevent"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	"errors"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"sync"
	"github.com/dgrijalva/jwt-go"
)

type Player struct {}

func (Player) Name(s string) string {
	return string((s[0])+32) + s[1:]
}

func (Player) PlayerSettingsGet(so *wsevent.Client, args struct {
	Key *string `json:"key"`
}) interface{} {

	player := chelpers.GetPlayer(so.Token)
	if *args.Key == "*" {
		return newResponse(player.Settings)
	}

	setting := player.GetSetting(*args.Key)
	return newResponse(setting)
}

func (Player) PlayerSettingsSet(so *wsevent.Client, args struct {
	Key *string `json:"key"`
	Value *string `json:"value"`
}) interface {} {

	player := chelpers.GetPlayer(so.Token)

	switch *args.Key {
	case "siteAlias":
		if len(*args.Value) > 32 {
			return errors.New("Site alias must be under 32 characters long.")
		}
		player.SetSetting(*args.Key, *args.Value)

		player.SetPlayerProfile()
		so.EmitJSON(helpers.NewRequest("playerProfile", player))

	default:
		player.SetSetting(*args.Key, *args.Value)
	}

	return emptySuccess
}

func (Player) PlayerProfile(so *wsevent.Client, args struct {
	Studentid *string `json:"studentid"`
}) interface {} {
	 studentid := *args.Studentid
	if studentid == "" {
		claims := so.Token.Claims.(jwt.MapClaims)
		studentid = claims["student_id"].(string)
	}

	player, err := player.GetPlayerByStudentID(studentid)
	if err != nil {
		return err
	}

	player.SetPlayerProfile()

	return newResponse(player)
}

var (
	changeMu = new(sync.RWMutex)
)