package player

import (
	"time"
	db "github.com/andreweggleston/GoSeniorAssassin/databaseAssassin"
)

type Report struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time

	PlayerID uint
	Type     ReportType
}

type ReportType int

const (
	Substitute ReportType = iota //!sub
	Vote                         //!repped by other players
	RageQuit                     //rage quit
)

func (player *Player) NewReport(rtype ReportType) {
	var count int

	last := time.Now().Add(-30 * time.Minute)
	db.DB.Model(&Report{}).Where("player_id = ? AND created_at > ? AND type = ?", player.ID, last, rtype).Count(&count)

	switch rtype {
	case Substitute:
		if count == 1 {
			player.BanUntil(time.Now().Add(30*time.Minute), BanJoin, "For !subbing twice in the last 30 minutes", 0)
		}
	case Vote:
		if count != 0 {
			player.BanUntil(time.Now().Add(30*time.Minute), BanJoin, "For getting reported multiple times", 0)
		}
	case RageQuit:
		if count != 0 {
			player.BanUntil(time.Now().Add(30*time.Minute), BanJoin, "For ragequitting a lobby multiple times in the last 30 minutes", 0)
		}

	}

	r := &Report{
		PlayerID: player.ID,
		Type:     rtype,
	}
	db.DB.Save(r)
}