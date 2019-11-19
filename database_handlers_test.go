package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var target = DownloadTarget{
	URL:                        "https://www.youtube.com/user/HungOverGargoyle",
	LatestDownloaded:           "123adadaf",
	DownloadMode:               "Audio Only",
	Name:                       "Electronic Gems",
	PreferredExtensionForAudio: "mp3",
	PreferredExtensionForVideo: "",
	DownloadHistory:            nil,
	LastChecked:                "",
	CheckingInterval:           "",
	Type:                       "Channel",
}

func TestDownloadTarget_AddToDatabase(t *testing.T) {
	assert := assert.New(t)
	err := target.AddToDatabase()
	assert.Nil(err)
	isTargetInDb, err := target.DoesExist()
	assert.Nil(err)
	assert.Equal(true, isTargetInDb, "isTargetInDb should be true")
	err = target.Delete()
	assert.Nil(err)
}

func TestDownloadTarget_UpdateLatestDownloaded(t *testing.T) {
	assert := assert.New(t)
	err := target.AddToDatabase()
	assert.Nil(err)
	err = target.UpdateLatestDownloaded("newvideoid123")
	assert.Nil(err)
	updatedTarget, err := target.GetFromDatabase()
	assert.Nil(err)
	assert.Equal("newvideoid123", updatedTarget.LatestDownloaded, "updatedTarget.LatestDownloaded should be newvideoid123")
	err = updatedTarget.Delete()
	assert.Nil(err)
}

func TestUpdateCheckingInterval(t *testing.T) {
	assert := assert.New(t)
	err := target.AddToDatabase()
	assert.Nil(err)
	_, err = UpdateCheckingInterval("5")
	assert.Nil(err)
	updatedTarget, err := target.GetFromDatabase()
	assert.Equal("5", updatedTarget.CheckingInterval, "updatedTarget.CheckingInterval should be 5")
	err = updatedTarget.Delete()
	assert.Nil(err)
}

func TestGetCheckingInterval(t *testing.T) {
	assert := assert.New(t)
	err := target.AddToDatabase()
	assert.Nil(err)
	_, err = UpdateCheckingInterval("5")
	assert.Nil(err)
	checkingInterval, err := GetCheckingInterval()
	assert.Nil(err)
	assert.Equal(5, checkingInterval, "checkingInterval should be 5")
	err = target.Delete()
	assert.Nil(err)
}

func TestDownloadTarget_UpdateLastChecked(t *testing.T) {
	currTime := time.Now().Format("01-02-2006 15:04:05")
	assert := assert.New(t)
	err := target.AddToDatabase()
	assert.Nil(err)
	err = target.UpdateLastChecked()
	assert.Nil(err)
	updatedTarget, err := target.GetFromDatabase()
	assert.Nil(err)
	assert.Equal(currTime, updatedTarget.LastChecked, "updatedTarget.LastChecked should be "+currTime)
	err = target.Delete()
	assert.Nil(err)
}