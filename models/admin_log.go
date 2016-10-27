package models

import (
	"github.com/jinzhu/gorm"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

type AdminLogEntry struct {
	gorm.Model
	PlayerID uint   //Admin responsible for action
	RelID    uint   `sql:"default:0"`  //The targated player
	RelText  string `sql:"default:''"` //The action text
}

func LogCustomAdminAction(playerid uint, reltext string, relid uint) error {
	entry := AdminLogEntry{
		PlayerID: playerid,
		RelID:    relid,
		RelText:  reltext,
	}

	return database.DB.Create(&entry).Error
}

func LogAdminAction(playerid uint, permission authority.AuthAction, relid uint) error {
	return LogCustomAdminAction(playerid, helpers.ActionNames[permission], relid)
}