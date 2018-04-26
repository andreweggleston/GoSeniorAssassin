package migrations

import (
	"github.com/sirupsen/logrus"
	db "github.com/andreweggleston/GoSeniorAssassin/databaseAssassin"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
)

var migrationRoutines = map[uint64]func(){
	3:  dropSubtituteTable,
	4:  increaseChatMessageLength,
	5:  updateAllPlayerInfo,
	6:  truncateHTTPSessions,
	7:  setMumbleInfo,
	8:  setPlayerExternalLinks,
	9:  setPlayerSettings,
	10: dropTableSessions,
	11: dropColumnUpdatedAt,
	12: moveReportsServers,
	13: dropUnusedColumns,
}





func dropSubtituteTable() {
	db.DB.Exec("DROP TABLE substitutes")
}

func increaseChatMessageLength() {
	db.DB.Exec("ALTER TABLE chat_messages ALTER COLUMN message TYPE character varying(150)")
}

func updateAllPlayerInfo() {
	var players []*player.Player
	db.DB.Model(&player.Player{}).Find(&players)

	for _, player := range players {
		player.Save()
	}
}

func truncateHTTPSessions() {
	db.DB.Exec("TRUNCATE TABLE http_sessions")
}

func setMumbleInfo() {
	var players []*player.Player

	db.DB.Model(&player.Player{}).Find(&players)
	for _, player := range players {
		player.Save()
	}
}

func setPlayerExternalLinks() {
	var players []*player.Player
	db.DB.Model(&player.Player{}).Find(&players)

	for _, player := range players {
		player.Save()
	}
}

// move player_settings values to player.Settings hstore
func setPlayerSettings() {
	rows, err := db.DB.DB().Query("SELECT player_id, key, value FROM player_settings")
	if err != nil {
		logrus.Fatal(err)
	}
	for rows.Next() {
		var playerID uint
		var key, value string

		rows.Scan(&playerID, &key, &value)
		p, _ := player.GetPlayerByID(playerID)
		p.SetSetting(key, value)
	}

	db.DB.Exec("DROP TABLE player_settings")
}

func dropTableSessions() {
	db.DB.Exec("DROP TABLE http_sessions")
}

func dropColumnUpdatedAt() {
	db.DB.Exec("ALTER TABLE players DROP COLUMN updated_at")
}

func moveReportsServers() {
	type oldReport struct {
		PlayerID uint
		LobbyID  uint
	}

	reportTypes := []struct {
		table string
		rtype player.ReportType
	}{
		{"ragequits_player_lobbies", player.RageQuit},
		{"reports_player_lobbies", player.Vote},
		{"substitutes_player_lobbies", player.Substitute},
	}

	for _, rtype := range reportTypes {
		rows, _ := db.DB.DB().Query("SELECT * FROM " + rtype.table)

		logrus.Info("Creating entries for ", rtype.table)
		for rows.Next() {
			var report oldReport

			rows.Scan(&report.PlayerID, &report.LobbyID)
			p, _ := player.GetPlayerByID(report.PlayerID)
			newReport := &player.Report{
				PlayerID: p.ID,
				Type:     rtype.rtype,
			}
			db.DB.Create(newReport)
		}
	}

	db.DB.Exec("DROP TABLE ragequits_player_lobbies")
	db.DB.Exec("DROP TABLE reports_player_lobbies")
	db.DB.Exec("DROP TABLE substitutes_player_lobbies")
}

func dropUnusedColumns() {
	db.DB.Model(&player.Player{}).DropColumn("debug")
}

