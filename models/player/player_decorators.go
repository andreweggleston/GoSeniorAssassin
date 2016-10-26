package player

import "github.com/andreweggleston/GoSeniorAssassin/helpers"

func (p *Player) DecoratePlayerTags() []string {
	tags := []string{helpers.RoleNames[p.Role]}
	return tags
}

func (p *Player) setJSONFields(bans bool) {


	p.PlaceholderTags = new([]string)
	p.PlaceholderRoleStr = new(string)

	*p.PlaceholderRoleStr = helpers.RoleNames[p.Role]
	*p.PlaceholderTags = p.DecoratePlayerTags()

	// if lobbies {
	// 	p.PlaceholderLobbies = new([]LobbyData)
	// 	var lobbies []*Lobby
	// 	db.DB.Table("lobbies").Joins("INNER JOIN lobby_slots ON lobbies.id = lobby_slots.lobby_id").Where("lobbies.match_ended = TRUE AND lobby_slots.player_id = ?", p.ID).Order("lobbies.ID DESC").Limit(5).Find(&lobbies)

	// 	for _, lobby := range lobbies {
	// 		*p.PlaceholderLobbies = append(*p.PlaceholderLobbies, DecorateLobbyData(lobby, true))
	// 	}
	// }

	p.Name = p.Alias()

	if bans {
		p.PlaceholderBans, _ = p.GetActiveBans()
	}
}

func (p *Player) SetPlayerProfile() {
	p.setJSONFields(true)
}

func (p *Player) SetPlayerSummary() {
	p.setJSONFields(false)
}