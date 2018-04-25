package database

import (
	"testing"
	"os"
	"github.com/sirupsen/logrus"
	"strconv"
	"github.com/stretchr/testify/assert"
	"github.com/andreweggleston/GoSeniorAssassin/config"
)

func TestDatabasePing(t *testing.T) {
	ci := os.Getenv("CI")
	if ci == "true" {
		config.Constants.DbUsername = "postgres"
		config.Constants.DbDatabase = "travis_ci_test"
		config.Constants.DbPassword = ""
	}

	logrus.Debug("[Test.Database] IsTest? " + strconv.FormatBool(IsTest))
	Init()
	assert.Nil(t, DB.DB().Ping())
}