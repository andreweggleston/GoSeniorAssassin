package controllers

import (
	"html/template"
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/models/player"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/internal/version"
	"fmt"
	"runtime"
	"github.com/sirupsen/logrus"
)

var (
	mainTempl *template.Template
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var p *player.Player
	token, err := controllerhelpers.GetToken(r)

	if err == nil {
		p = controllerhelpers.GetPlayer(token)
	}

	errtempl := mainTempl.Execute(w, map[string]interface{}{
		"LoggedIn": 	err == nil,
		"Player":	p,
		"MockLogin": 	config.Constants.MockupAuth,
		"BuildDate": 	version.BuildDate,
		"GitCommit":	version.GitCommit,
		"GitBranch":	version.GitBranch,
		"BuildInfo": fmt.Sprintf("Build using %s on %s (%s %s)", runtime.Version(), version.Hostname, runtime.GOOS, runtime.GOARCH),
	})
	if errtempl != nil {
		logrus.Error(err)
	}
}
