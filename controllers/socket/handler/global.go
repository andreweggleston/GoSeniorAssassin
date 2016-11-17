package handler

import (
	"github.com/TF2Stadium/wsevent"
	"github.com/bitly/go-simplejson"
	"errors"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/broadcaster"
	"github.com/dgrijalva/jwt-go"
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
	claims := so.Token.Claims.(jwt.MapClaims)
	ID := claims["id"].(string)
	broadcaster.SendMessageSkipIDs(so.ID, ID, args.Event, args.Data)
	return emptySuccess
}
