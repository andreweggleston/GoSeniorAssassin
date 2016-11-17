package controllerhelpers

import (
	"sync"
	"time"
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/Sirupsen/logrus"
	"encoding/xml"
	"github.com/andreweggleston/GoSeniorAssassin/helpers/authority"
	"github.com/TF2Stadium/wsevent"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"errors"
	"strconv"
	"github.com/dgrijalva/jwt-go"
)

var (
	whitelistLock = new(sync.RWMutex)
	whitelistID map[string]bool
)

func WhitelistListener() {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(config.Constants.IDWhitelist)

		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		var groupXML struct {
			Members []string `xml:"members>ID64"`
		}

		dec := xml.NewDecoder(resp.Body)
		err = dec.Decode(&groupXML)
		if err != nil {
			logrus.Error(err)
			continue
		}

		whitelistLock.Lock()
		whitelistID = make(map[string]bool)

		for _, ID := range groupXML.Members {
			whitelistID[ID] = true
		}
		whitelistLock.Unlock()
		<-ticker.C
	}
}
func IsIDWhitelisted(id string) bool {
	whitelistLock.RLock()
	defer whitelistLock.RUnlock()
	whitelisted, exists := whitelistID[id]


	return whitelisted && exists
}

func CheckPrivilege(so *wsevent.Client, action authority.AuthAction) error {
	claims := so.Token.Claims.(jwt.MapClaims)
	player, _ := player.GetPlayerByID(claims["id"].(uint))
	if !player.Role.Can(action) {
		return errors.New("You are not authorized to perform this action")
	}
	return nil
}

func FilterHTTPRequest(action authority.AuthAction, f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetToken(r)
		if err != nil {
			http.Error(w, "you're not logged in, or your JWT cookie is invalid", http.StatusBadRequest)
			return
		}

		claims:=token.Claims.(jwt.MapClaims)
		role, _ := strconv.Atoi(claims["role"].(string))

		if !(authority.AuthRole(role).Can(action)) {
			http.Error(w, "Not authorized", 403)
			return
		}

		f(w, r)
	}
}