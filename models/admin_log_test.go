package models

import (
	"github.com/andreweggleston/GoSeniorAssassin/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/andreweggleston/GoSeniorAssassin/database"
	"github.com/andreweggleston/GoSeniorAssassin/helpers"
)

func init() {
	testhelpers.CleanupDB()
}

func TestLogCreation(t *testing.T) {
	t.Parallel()
	var obj = AdminLogEntry{}
	count := 5
	database.DB.Model(obj).Count(&count)
	assert.Equal(t, 0, count)

	LogAdminAction(1, helpers.ActionBanJoin, 2)
	LogCustomAdminAction(2, "test", 4)

	database.DB.Model(obj).Count(&count)
	assert.Equal(t, 2, count)
}