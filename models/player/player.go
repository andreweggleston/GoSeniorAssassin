package player

import (

)
import (
	"errors"
	"time"
	"github.com/jinzhu/gorm/dialects/postgres"
	db "github.com/drewwww/SeniorAssassin/database"
)

var ErrPlayerNotFound = errors.New("Player not found")
var ErrPlayerInReportedSlot = errors.New("Player in reported slot")

type Player struct {
	ID			uint		`gorm:"primary_key" json:"id"`
	CreatedAt            	time.Time 	`json:"createdAt"`
	ProfileUpdatedAt      	time.Time 	`json:"-"`
	StreamStatusUpdatedAt 	time.Time 	`json:"-"`

	Settings postgres.Hstore `json:"-"`
}

func NewPlayer(Id string) (*Player, error) {
	player := &Player{ID: Id}

	last := &Player{}
	db.DB.Model(&Player{}).Last(last)

	return player, nil
}
