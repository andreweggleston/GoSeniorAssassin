package controllerhelpers

import (
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/Sirupsen/logrus"
	"encoding/base64"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
	"fmt"
	"net/http"
)

var (
	signingKey []byte
)

func init() {
	if config.Constants.CookieStoreSecret == "secret" {
		logrus.Warning("Using an insecure encruption key")
		signingKey = []byte("secret")
		return
	}

	var err error
	signingKey, err = base64.StdEncoding.DecodeString(config.Constants.CookieStoreSecret)
	if err != nil {
		logrus.Fatal(err)
	}
}

func NewToken(player *player.Player) string {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := make(jwt.MapClaims)
	claims["player_id"] = strconv.FormatUint(uint64(player.ID), 10)
	claims["role"] = strconv.Itoa(int(player.Role))
	claims["iat"] = time.Now().Unix()
	claims["iss"] = config.Constants.PublicAddress
	token.Claims = claims

	str, err := token.SignedString([]byte(signingKey))
	if err != nil {
		logrus.Error(err)
	}

	return str
}

func verifyToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return signingKey, nil
}

func GetToken(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("auth-jwt")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(cookie.Value, verifyToken)
	return token, err
}

func GetPlayer(token *jwt.Token) *player.Player {
	claims:=token.Claims.(jwt.MapClaims)
	playerid, _ := strconv.ParseUint(claims["player_id"].(string), 10, 32)
	player, _ := player.GetPlayerByID(uint(playerid))
	return player
}
