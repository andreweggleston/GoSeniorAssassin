package database

import (
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/andreweggleston/GoSeniorAssassin/config"
)

var (
	IsTest		bool = false
	DB		*gorm.DB
	dbMutex		sync.Mutex
	initialized	= false
	DBUrl		url.URL
)

func Init() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if initialized{
		return
	}

	DBUrl = url.URL{
		Scheme:	"postgres",
		Host:	config.Constants.DbAddr,
		RawQuery: "sslmode=disable",
	}

	logrus.Info("Connecting to DB on ", DBUrl.String())

	DBUrl.User = url.UserPassword(config.Constants.DbUsername, config.Constants.DbPassword)

	var err error
	DB, err = gorm.Open("postgres", DBUrl.String())
	if err != nil {
		logrus.Fatal(err.Error())
	}

	DB.SetLogger(logrus.StandardLogger())

	logrus.Info("Connected!")
	initialized = true
}
