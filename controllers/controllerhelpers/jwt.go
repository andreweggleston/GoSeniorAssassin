package controllerhelpers

import (
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/sirupsen/logrus"
	"encoding/base64"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"net/http"
)

var (
	signingKey []byte
)

func init() {
	if config.Constants.CookieStoreSecret == "secret" {
		logrus.Warning("Using an insecure encryption key")
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
	token.Claims = AssassinClaims{
		PlayerID:  		player.ID,
		StudentID:		player.StudentID,
		Role: 			player.Role,
		IssuedAt: 		time.Now().Unix(),
		Issuer: 		config.Constants.PublicAddress,

	}

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

	token, err := jwt.ParseWithClaims(cookie.Value, &AssassinClaims{},verifyToken)
	return token, err
}

func GetPlayer(token *jwt.Token) *player.Player {

	player, _ := player.GetPlayerByID(token.Claims.(*AssassinClaims).PlayerID)
	return player
}
