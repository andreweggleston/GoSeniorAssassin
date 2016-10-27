package testhelpers

import (
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

func CreatePlayer() *player.Player {
	studentID := string("andrew_eggleston")

	player, _ := player.NewPlayer(studentID)
	player.Save()
	return player
}

func CreatePlayerMod() *player.Player {
	p := CreatePlayer()
	p.Role = helpers.RoleMod
	p.Save()
	return p
}

func CreatePlayerAdmin() *player.Player {
	p := CreatePlayer()
	p.Role = helpers.RoleAdmin
	p.Save()
	return p
}