package controllerhelpers

import (
	"github.com/sfreiberg/gotwilio"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)

var (
	accountSid = "ACd859873daf19b011ffb8ce9876450a40"
	authToken  = "15049643a7e4ee4ae752cdedae625cf4"
	twilio     = gotwilio.NewTwilioClient(accountSid, authToken)
	from       = "+15083881310"

	playerArray	[]player.Player
)

func SendGlobalMessage(message string) {

	playerlist, _ := database.DB.Not("phone_number", "").Find(&playerArray).Rows()

	defer playerlist.Close()
	for playerlist.Next() {
		var player player.Player
		database.DB.ScanRows(playerlist, &player)

		twilio.SendSMS(from, "+1"+player.PhoneNumber, message, "", "")
	}

}

func SendMessage(phoneNumber string, message string){
	twilio.SendSMS(from, "+1"+phoneNumber, message, "", "")
}