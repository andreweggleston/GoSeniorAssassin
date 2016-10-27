package chat

import (
	"fmt"
	"time"

	"encoding/json"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/broadcaster"
	db "github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)


// ChatMessage Represents a chat mesasge sent by a particular player
type ChatMessage struct {
									 // Message ID
	ID        uint          `json:"id"`
	CreatedAt time.Time     `json:"timestamp"`
	Player    player.Player `json:"-"`
	PlayerID  uint          `json:"-"`                               // ID of the player who sent the message
	Room      int           `json:"room"`                            // room to which the message was sent
	Message   string        `json:"message" sql:"type:varchar(150)"` // the actual Message
	Deleted   bool          `json:"deleted"`                         // true if the message has been deleted by a moderator
	Bot       bool          `json:"bot"`                             // true if the message was sent by the notification "bot"
}

// Return a new ChatMessage sent from specficied player
func NewChatMessage(message string, room int, player *player.Player) *ChatMessage {
	record := &ChatMessage{
		PlayerID: player.ID,

		Room:    room,
		Message: message,
	}

	return record
}


func (m *ChatMessage) Save() {
	db.DB.Save(m)
}

func (m *ChatMessage) Send() {
	broadcaster.SendMessageToRoom(fmt.Sprintf("%d_public", m.Room), "chatReceive", m)
	if m.Room != 0 {
		broadcaster.SendMessageToRoom(fmt.Sprintf("%d_private", m.Room), "chatReceive", m)
	}
}

// we only need these three things for showing player messages
type minPlayer struct {
	Name    string   `json:"name"`
	StudentID string   `json:"studentid"`
	Tags    []string `json:"tags"`
}

var bot = minPlayer{"AssassinBot", "assassin_bot", []string{"tf2stadium"}}

func (m *ChatMessage) MarshalJSON() ([]byte, error) {
	message := map[string]interface{}{
		"id":        m.ID,
		"timestamp": m.CreatedAt,
		"room":      m.Room,
		"message":   m.Message,
		"deleted":   m.Deleted,
	}
	if m.Bot {
		message["player"] = bot
	} else {
		p := &player.Player{}
		db.DB.First(p, m.PlayerID)
		player := minPlayer{
			Name:	p.StudentID,
			Tags:	p.DecoratePlayerTags(),
		}

		if m.Deleted {
			player.Tags = append(player.Tags, "<deleted>")
			message["message"] = "<deleted>"
		}

		message["player"] = player
	}


	return json.Marshal(message)
}

func NewBotMessage(message string, room int) *ChatMessage {
	m := &ChatMessage{
		Room:    room,
		Message: message,

		Bot: true,
	}

	m.Save()
	return m
}

func SendNotification(message string, room int) {
	pub := fmt.Sprintf("%d_public", room)
	broadcaster.SendMessageToRoom(pub, "chatReceive", NewBotMessage(message, room))
}

// Return a list of ChatMessages spoken in room
func GetRoomMessages(room int) ([]*ChatMessage, error) {
	var messages []*ChatMessage

	err := db.DB.Model(&ChatMessage{}).Where("room = ?", room).Order("created_at").Find(&messages).Error

	return messages, err
}

// Return all messages sent by player to room
func GetPlayerMessages(p *player.Player) ([]*ChatMessage, error) {
	var messages []*ChatMessage

	err := db.DB.Model(&ChatMessage{}).Where("player_id = ?", p.ID).Order("room, created_at").Find(&messages).Error

	return messages, err

}

// Get a list of last 20 messages sent to room, used by frontend for displaying the chat history/scrollback
func GetScrollback(room int) ([]*ChatMessage, error) {
	var messages []*ChatMessage // apparently the ORM works fine with using this type (they're aliases after all)

	err := db.DB.Table("chat_messages").Where("room = ? AND deleted = FALSE", room).Order("id desc").Limit(20).Find(&messages).Error

	return messages, err
}
