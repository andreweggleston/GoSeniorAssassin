package handler

import (
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	db "github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/TF2Stadium/wsevent"
	"fmt"
	"time"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"errors"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
	"strings"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

type Chat struct{}

func (Chat) Name(s string) string {
	return string((s[0]) + 32) + s[1:]
}

func (Chat) ChatSend(so *wsevent.Client, args struct {
	Message *string `json:"message"`
	Room    *int        `json:"room"`
}) interface{} {
	p := chelpers.GetPlayer(so.Token)
	if banned, until := p.IsBannedWithTime(player.BanChat); banned {
		ban, _ := p.GetActiveBan(player.BanChat)
		return fmt.Errorf("You've been banned until %s (%s)", until.Format(time.RFC822), ban.Reason)
	}

	switch {
	case len(*args.Message) == 0:
		return errors.New("Cannot send an empty message")

	case (*args.Message)[0] == '\n':
		return errors.New("Cannot send messages prefixed with newline")

	case len(*args.Message) > 150:
		return errors.New("Message too long")
	}

	message := chat.NewChatMessage(*args.Message, *args.Room, p)

	if strings.HasPrefix(*args.Message, "!admin") {
		chelpers.SendToSlack(*args.Message, p.Name, p.StudentID)
		return emptySuccess
	}

	message.Save()
	message.Send()

	return emptySuccess
}

func (Chat) ChatDelete(so *wsevent.Client, args struct{
	ID *int `json:"id"`
	Room *uint `json:"room"`
}) interface{} {

	if err := chelpers.CheckPrivilege(so, helpers.ActionDeleteChat); err != nil {
		return err
	}

	message := &chat.ChatMessage{}
	err := db.DB.First(message, *args.ID).Error
	if message.Bot {
		return errors.New("Cannot delete notification messages")
	}
	if err != nil {
		return errors.New("Can't find message")
	}

	message.Deleted = true
	message.Save()
	message.Send()

	return emptySuccess
}