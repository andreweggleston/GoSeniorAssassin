package handler

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/bitly/go-simplejson"
	"errors"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/broadcaster"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
)

type Global struct{}

func (Global) Name(s string) string {
	return string((s[0])+32) + s[1:]
}

func (Global) GetConstant(so *wsevent.Client, args struct {
	Constant string `json:"constant"`
}) interface{} {

	output := simplejson.New()
	switch args.Constant {
	default:
		return errors.New("Unkown constant.")
	}

	return newResponse(output)
}

func (Global) SendToOtherClients(so *wsevent.Client, args struct{
	Event string `json:"event"`
	Data string `json:"data"`
}) interface{} {
	ID := so.Token.Claims.(*chelpers.AssassinClaims).StudentID
	broadcaster.SendMessageSkipIDs(so.ID, ID, args.Event, args.Data)
	return emptySuccess
}
