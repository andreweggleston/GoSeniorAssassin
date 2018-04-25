package controllerhelpers

import (
	db "github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)

type AssassinClaims struct {
	PlayerID       uint               `json:"player_id"`
	StudentID      string             `json:"student_id"`
	MumblePassword string             `json:"mumble_password"`
	Role           authority.AuthRole `json:"role"`
	IssuedAt       int64              `json:"iat"`
	Issuer         string             `json:"iss"`
}

func playerExists(id uint, studentID string) bool {
	var count int
	db.DB.Model(&player.Player{}).Where("id = ? AND student_id = ?", id, studentID).Count(&count)
	return count != 0
}

func (c AssassinClaims) Valid() error {
	if !playerExists(c.PlayerID, c.StudentID) {
		return player.ErrPlayerNotFound
	}

	return nil
}
