package chat

import (
	"github.com/andreweggleston/GoSeniorAssassin/inside/testhelpers"
	"testing"
	"strconv"
	"github.com/stretchr/testify/assert"
	db "github.com/andreweggleston/GoSeniorAssassin/database"
)

func init() {
	testhelpers.CleanupDB()
}

func TestNewChatMessage(t *testing.T) {
	player := testhelpers.CreatePlayer()
	player.Save()

	for i := 0; i < 3; i++ {
		message := NewChatMessage(strconv.Itoa(i), 0, player)
		assert.NotNil(t, message)

		err := db.DB.Save(message).Error
		assert.Nil(t, err)
	}

	messages, err := GetRoomMessages(0)
	assert.Nil(t, err)
	assert.Equal(t, len(messages), 3)
}
