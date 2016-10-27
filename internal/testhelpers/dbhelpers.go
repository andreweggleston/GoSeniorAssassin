package testhelpers

import (
	"os"
	"sync"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/database/migrations"
)

var cleaningMutex sync.Mutex

var o = new(sync.Once)

func CleanupDB() {
	cleaningMutex.Lock()
	defer cleaningMutex.Unlock()

	o.Do(func() {
		ci := os.Getenv("CI")
		config.Constants.DbAddr = "127.0.0.1:5432"

		if ci == "true" {
			config.Constants.DbUsername = "postgres"
			config.Constants.DbDatabase = "travis_ci_test"
			config.Constants.DbPassword = ""
		} else {
			config.Constants.DbDatabase = "TESTseniorassassin"
			config.Constants.DbUsername = "TESTseniorassassin"
			config.Constants.DbPassword = "assassinpass"
		}

		database.Init()
		migrations.Do()
	})

	tables := []string{
		"admin_log_entries",
		"chat_messages",
		"player_bans",
		"players",
		"reports",
	}
	for _, table := range tables {
		database.DB.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY")
	}

}
