package migrations

import (
	"sync"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/models"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
)

var once = new(sync.Once)

func Do() {
	database.DB.Exec("CREATE EXTENSION IF NOT EXISTS hstore")
	database.DB.AutoMigrate(&player.Player{})
	database.DB.AutoMigrate(&models.AdminLogEntry{})
	database.DB.AutoMigrate(&player.PlayerBan{})
	database.DB.AutoMigrate(&chat.ChatMessage{})
	database.DB.AutoMigrate(&Constant{})
	database.DB.AutoMigrate(&player.Report{})


	once.Do(checkSchema)
}