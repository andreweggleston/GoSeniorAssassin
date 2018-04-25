package admin

import (
	"regexp"
	"html/template"
	db "github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/jinzhu/gorm"
	"net/http"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
)

var (
	dateRegex	= regexp.MustCompile(`(\d{2})-(\d{2})-(\d{4})`)
	chatLogsTempl *template.Template
)

func getPlayerID(ID string) (playerID uint) {
	db.DB.Model(&player.Player{}).Select("id").Where("id = ?", ID).Row().Scan(&playerID)
	return
}

func GetChatLogs(w http.ResponseWriter, r *http.Request) {
	var messages []*chat.ChatMessage
	values := r.URL.Query()

	room, err := strconv.Atoi(values.Get("room"))
	if err != nil && values.Get("room") != "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ID := values.Get("id")
	var from, to time.Time

	if values.Get("from") != "" { //2006-01-02
		from, err = time.Parse("2006-01-02", values.Get("from"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		from = time.Time{}
	}

	if values.Get("to") != "" {
		to, err = time.Parse("2006-01-02", values.Get("to"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		to = time.Now()
	}

	order := values.Get("order")
	var results *gorm.DB

	if values.Get("room") == "" { //Retrieve all messages sent by a specific player
		if ID == "" {
			http.Error(w, "No ID given.", http.StatusBadRequest)
			return
		}

		playerID := getPlayerID(ID)
		if playerID == 0 {
			http.Error(w, fmt.Sprintf("Couldn't find player with ID %s", ID), http.StatusNotFound)
			return
		}

		results = db.DB.Preload("Player").Where("player_id = ? AND room = ? AND created_at >= ? AND created_at <= ?", playerID, room, from, to)
	} else if ID == "" { //Retrieve all messages sent to a specfic room
		results = db.DB.Preload("Player").Where("room = ? AND (created_at >= ? AND created_at <= ?)", room, from, to)
	} else { //Retrieve all messages sent to a specific room and a speficic player
		playerID := getPlayerID(ID)
		if playerID == 0 {
			http.Error(w, fmt.Sprintf("Couldn't find player with ID %s", ID), http.StatusNotFound)
			return
		}

		results = db.DB.Preload("Player").Where("player_id = ? AND room = ? AND created_at >= ? AND created_at <= ?", playerID, room, from, to)
	}

	if order == "Ascending" {
		err = results.Order("id").Find(&messages).Error
	} else if order == "Descending" {
		err = results.Order("id desc").Find(&messages).Error
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = chatLogsTempl.Execute(w, messages)
	if err != nil {
		logrus.Error(err)
	}
}