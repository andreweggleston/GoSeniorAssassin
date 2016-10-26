package player

import (

)
import (
	"errors"
	"time"
	"github.com/jinzhu/gorm/dialects/postgres"
	db "github.com/drewwww/SeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
)

var ErrPlayerNotFound = errors.New("Player not found")
var ErrPlayerInReportedSlot = errors.New("Player in reported slot")

type Player struct {
	ID			uint		`gorm:"primary_key" json:"id"`
	Name			string		`json:"name"`
	CreatedAt            	time.Time 	`json:"createdAt"`
	ProfileUpdatedAt      	time.Time 	`json:"-"`
	StreamStatusUpdatedAt 	time.Time 	`json:"-"`

	Settings postgres.Hstore `json:"-"`

	Role       authority.AuthRole `sql:"default:0" json:"-"`

	JSONFields
}

type JSONFields struct {
	PlaceholderTags          *[]string `sql:"-" json:"tags"`
	PlaceholderRoleStr       *string   `sql:"-" json:"role"`
	PlaceholderBans  []*PlayerBan `sql:"-" json:"bans"`
}

func NewPlayer(Id string) (*Player, error) {
	player := &Player{ID: Id}

	last := &Player{}
	db.DB.Model(&Player{}).Last(last)

	return player, nil
}

func isClean(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, c := range s {
		if c >= 'A' && c <= 'z' || c == ' ' {
			continue
		}
		return false
	}

	return true
}

func (p *Player) Alias() string {
	alias := p.GetSetting("siteAlias")
	if alias == "" {
		return p.Name
	}

	return alias
}

func (player *Player) Save() error {
	var err error
	if db.DB.NewRecord(player) {
		err = db.DB.Create(player).Error
	} else {
		err = db.DB.Save(player).Error
	}
	return err
}

func GetPlayerByID(ID uint) (*Player, error) {
	player := &Player{}

	if err := db.DB.First(player, ID).Error; err!= nil {
		return nil, err
	}

	return player, nil
}

func (player *Player) GetSetting(key string) string {
	if player.Settings == nil {
		return ""
	}

	value, ok := player.Settings[key]
	if !ok {
		return ""
	}

	return *value
}
func (player *Player) SetSetting(key string, value string) {
	if player.Settings == nil {
		player.Settings = make(postgres.Hstore)
	}

	player.Settings[key] = &value
	player.Save()
}

func (player *Player) UpdatePlayerInfo() error {
	//TODO: get elbing to make some shitty thing to grab data from ipass

	return nil
}
