package migrations

import (
	"sync"
	"github.com/andreweggleston/GoSeniorAssassin/databaseAssassin"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/models"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
)

var once = new(sync.Once)

func Do() {
	databaseAssassin.DB.Exec("CREATE EXTENSION IF NOT EXISTS hstore")
	databaseAssassin.DB.AutoMigrate(&player.Player{})
	databaseAssassin.DB.AutoMigrate(&models.AdminLogEntry{})
	databaseAssassin.DB.AutoMigrate(&player.PlayerBan{})
	databaseAssassin.DB.AutoMigrate(&chat.ChatMessage{})
	databaseAssassin.DB.AutoMigrate(&Constant{})
	databaseAssassin.DB.AutoMigrate(&player.Report{})


	once.Do(checkSchema)
}