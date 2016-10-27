package player

import (
	"github.com/andreweggleston/GoSeniorAssassin/internal/testhelpers"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func init() {
	testhelpers.CleanupDB()
}

func TestGetPlayerByStudentID(t *testing.T) {
	t.Parallel()
	player := testhelpers.CreatePlayer()
	player2, err := GetPlayerByID(player.ID)
	assert.NoError(t, err)
	assert.Equal(t, player.ID, player2.ID)
}

func TestPlayerSettings(t *testing.T) {
	t.Parallel()

	player := testhelpers.CreatePlayer()

	settings := player.Settings
	assert.Equal(t, 0, len(settings))

	player.SetSetting("foo", "bar")
	assert.Equal(t, player.GetSetting("foo"), "bar")

	player.SetSetting("hello", "world")
	assert.Equal(t, player.GetSetting("hello"), "world")
	assert.Len(t, player.Settings, 2)
}

func TestPlayerBanning(t *testing.T) {
	t.Parallel()
	player := testhelpers.CreatePlayer()

	for ban := BanJoin; ban != BanFull; ban++ {
		assert.False(t, player.IsBanned(ban))
	}

	past := time.Now().Add(time.Second * -10)
	player.BanUntil(past, BanJoin, "they suck", 0)
	assert.False(t, player.IsBanned(BanJoin))

	future := time.Now().Add(time.Second * 10)
	player.BanUntil(future, BanJoin, "they suck", 0)
	player.BanUntil(future, BanFull, "they suck", 0)

	player2, _ := GetPlayerByStudentID(player.StudentID)
	assert.True(t, player2.IsBanned(BanChat))
	isBannedFull, untilFull := player2.IsBannedWithTime(BanFull)
	assert.True(t, isBannedFull)
	assert.True(t, future.Sub(untilFull) < time.Second)
	assert.True(t, untilFull.Sub(future) < time.Second)


	future2 := time.Now().Add(time.Second * 20)
	player2.BanUntil(future2, BanJoin, "they suck", 0)

	bans, err := player2.GetActiveBans()
	assert.NoError(t, err)
	assert.Len(t, bans, 2)

	_, err = player2.GetActiveBan(BanJoin)
	assert.NoError(t, err)

	player2.Unban(BanJoin)
	player2.Unban(BanFull)

	for ban := BanJoin; ban != BanFull; ban++ {
		assert.False(t, player2.IsBanned(ban))
	}
}