package player

import (
	"errors"
	"time"
	"github.com/jinzhu/gorm/dialects/postgres"
	db "github.com/andreweggleston/GoSeniorAssassin/databaseAssassin"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"strings"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	"github.com/sirupsen/logrus"
	"github.com/sfreiberg/gotwilio"
)

var ErrPlayerNotFound = errors.New("Player not found")
var ErrPlayerInReportedSlot = errors.New("Player in reported slot")

type Player struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Sub         string `sql:"not null;unique" json:"sub"`
	StudentID   string `sql:"not null;unique" json:"studentid"`
	Name        string `json:"name"`
	Email       string `sql:"not null" json:"email"`
	PhoneNumber string `json:"phonenumber"`

	Target           string    `json:"target"`
	Kills            uint      `sql:"not null" json:"kills"`
	CreatedAt        time.Time `json:"createdAt"`
	ProfileUpdatedAt time.Time `json:"-"`

	MarkedForDeath bool `sql:"not null" json:"markedfordeath"`
	Killed         bool `sql:"not null" json:"killed"`

	GlobalData

	Settings postgres.Hstore `json:"-"`

	Role authority.AuthRole `sql:"default:0" json:"-"`

	JSONFields
}

type JSONFields struct {
	PlaceholderTags    *[]string    `sql:"-" json:"tags"`
	PlaceholderRoleStr *string      `sql:"-" json:"role"`
	PlaceholderBans    []*PlayerBan `sql:"-" json:"bans"`
}

type GlobalData struct {
	SafetyItem   string
	KillByDate   string
	Announcement string
}

func NewPlayer(studentID string) (*Player, error) {
	player := &Player{StudentID: studentID, Kills: 0}

	if isAdmin(studentID) {
		player.Role = helpers.RoleDeveloper
	}

	last := &Player{}
	db.DB.Model(&Player{}).Last(last)

	return player, nil
}

func isAdmin(studentID string) bool {

	for _, b := range strings.Split(config.Constants.AdminStudentIDs, ",") {
		if studentID == b {
			return true
		}
	}
	return false
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

	if err := db.DB.First(player, ID).Error; err != nil {
		return nil, err
	}

	return player, nil
}

func GetPlayerByStudentID(studentid string) (*Player, error) {
	var player = Player{}
	err := db.DB.Where("student_id = ?", studentid).First(&player).Error
	if err != nil {
		return nil, ErrPlayerNotFound
	}
	return &player, nil
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

	if key == "phoneNumber" {
		player.PhoneNumber = value
	}

	player.Settings[key] = &value
	player.Save()
}

func (player *Player) MarkForDeath() {
	player.MarkedForDeath = true
	player.Save()
	accountSid := "ACd859873daf19b011ffb8ce9876450a40"
	authToken  := "15049643a7e4ee4ae752cdedae625cf4"
	twilio     := gotwilio.NewTwilioClient(accountSid, authToken)
	twilio.SendSMS("+15083881310", "+1"+player.PhoneNumber, "You've been marked for death. Have you been assassinated? Confirm or deny on the website.", "", "")
}

func (player *Player) MarkTarget() error {
	target, err := GetPlayerByStudentID(player.Target)

	target.MarkForDeath()

	return err
}

func (player *Player) ConfirmOwnMark() {
	if player.MarkedForDeath {
		player.Killed = true
		player.Save()
	}
}

func (player *Player) DenyOwnMark() {
	if player.MarkedForDeath {
		player.MarkedForDeath = false
		player.Save()
	}
}

func (player *Player) TargetIsDead() bool {
	if player.Target != "" {
		target, err := GetPlayerByStudentID(player.Target)
		if err != nil {
			logrus.Error(err)
		}
		return target.Killed
	}
	return false

}
func (player *Player) UpdatePlayerData() error {
	defer player.Save()
	if player.TargetIsDead() {
		target, err := GetPlayerByStudentID(player.Target)
		if err != nil {
			return err
		}
		player.Target = target.Target
		player.Kills = player.Kills + 1

		accountSid := "ACd859873daf19b011ffb8ce9876450a40"
		authToken  := "15049643a7e4ee4ae752cdedae625cf4"
		twilio     := gotwilio.NewTwilioClient(accountSid, authToken)
		twilio.SendSMS("+15083881310", "+1"+player.PhoneNumber, "Your target has been updated.", "", "")

	}
	return nil
}
