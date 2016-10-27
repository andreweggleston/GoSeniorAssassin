package main

import (
	"flag"
	"crypto/rand"
	"github.com/Sirupsen/logrus"
	"encoding/base64"
	"fmt"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"os"
	"github.com/andreweggleston/GoSeniorAssassin/internal/version"
	"github.com/andreweggleston/GoSeniorAssassin/controllers"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/database/migrations"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
	chelpers "github.com/andreweggleston/GoSeniorAssassin/controllers/controllerhelpers"
	"net/http"
	"github.com/andreweggleston/GoSeniorAssassin/routes"
	socketServer "github.com/andreweggleston/GoSeniorAssassin/routes/socket"
	"github.com/andreweggleston/GoSeniorAssassin/controllers/socket"
	"github.com/rs/cors"
	"os/signal"
	"syscall"
	"github.com/andreweggleston/GoSeniorAssassin/models/chat"
)

var (
	flagGen = flag.Bool("genkey", false, "write a 32bit key for encrypting cookies the given file, and exit")
	docPrint = flag.Bool("printdoc", false, "print the docs for environment variables, and exit.")
	dbMaxopen = flag.Int("db-maxopen", 80, "maximum number of open database connections allowed.")
)

func main() {
	flag.Parse()

	if *flagGen {
		key := make([]byte, 64)
		_, err := rand.Read(key)
		if err != nil {
			logrus.Fatal(err)
		}

		base64Key := base64.StdEncoding.EncodeToString(key)
		fmt.Println(base64Key)
		return
	}
	if *docPrint {
		config.PrintConfigDoc()
		os.Exit(0)
	}

	logrus.Debug("Commit: ", version.GitCommit)
	logrus.Debug("Branch: ", version.GitBranch)
	logrus.Debug("Build date: ", version.BuildDate)

	controllers.InitTemplates()

	database.Init()
	database.DB.DB().SetMaxOpenConns(*dbMaxopen)
	migrations.Do()

	helpers.ConnectAMQP()

	if config.Constants.IDWhitelist != "" {
		go chelpers.WhitelistListener()
	}

	httpMux := http.NewServeMux()
	routes.SetupHTTP(httpMux)
	socket.RegisterHandlers()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:		config.Constants.AllowedOrigins,
		AllowedMethods:		[]string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowCredentials:	true,
	}).Handler(httpMux)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-sig
		shutdown()
		os.Exit(0)
	}()

	logrus.Info("Serving on ", config.Constants.ListenAddress)
	logrus.Info("Hosting on ", config.Constants.PublicAddress)

	logrus.Fatal(http.ListenAndServe(config.Constants.ListenAddress, corsHandler))
}

func shutdown() {
	logrus.Info("Received SIGINT/SIGTERM")
	chat.SendNotification(`Backend will be going down for a while for an update, click on "Reconnect" to reconnect to Senior Assassin`, 0)
	logrus.Info("waiting for GlobalWait")
	helpers.GlobalWait.Wait()
	logrus.Info("waiting for socket requests to complete.")
	socketServer.Wait()
	logrus.Info("closing all active websocket connections")
	socketServer.AuthServer.Close()
}
