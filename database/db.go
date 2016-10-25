package database

import (
	"net/url"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/drewwww/owbackend/config"
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
		Path:	config.Constants.DbDatabase,
		RawQuery: "sslmode=disable",
	}

	logrus.Info("Connecting to DB on ", DBUrl.String())
}