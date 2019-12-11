package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var targets = []DownloadTarget{
	{
		URL:                        "https://www.youtube.com/user/HungOverGargoyle",
		LatestDownloaded:           "electronicgemslatestdownloaded123",
		DownloadMode:               "Audio Only",
		Name:                       "Electronic Gems",
		PreferredExtensionForAudio: "mp3",
		PreferredExtensionForVideo: "",
		DownloadHistory:            nil,
		LastChecked:                "",
		CheckingInterval:           "",
		Type:                       "Channel",
	},
	{
		URL:                        "https://www.youtube.com/user/NewRetroWave",
		LatestDownloaded:           "newretrowavelatestdownloaded123",
		DownloadMode:               "Video And Audio",
		Name:                       "NewRetroWave",
		PreferredExtensionForAudio: "",
		PreferredExtensionForVideo: "mp4",
		DownloadHistory:            nil,
		LastChecked:                "",
		CheckingInterval:           "",
		Type:                       "Channel",
	},
	{
		URL:                        "https://www.youtube.com/watch?v=e0ML7i24trI&list=PLyIFQr1wryPJvYmJ5JDtRE5tH78UhyB7s",
		LatestDownloaded:           "newretrowavelatestdownloaded123",
		DownloadMode:               "Video And Audio",
		Name:                       "Heroic Motivational Anthems",
		PreferredExtensionForAudio: "",
		PreferredExtensionForVideo: "mp4",
		DownloadHistory:            nil,
		LastChecked:                "",
		CheckingInterval:           "5",
		Type:                       "Playlist",
	},
}

func TestDownloadTarget_AddToDatabase(t *testing.T) {
	assert := assert.New(t)
	for _, target := range targets {
		err := target.AddToDatabase()
		assert.Nil(err)
		isTargetInDb, err := target.DoesExist()
		assert.Nil(err)
		assert.Equal(true, isTargetInDb, "isTargetInDb should be true")
		err = target.Delete()
		assert.Nil(err)
	}
}

//
//func TestDownloadTarget_UpdateLatestDownloaded(t *testing.T) {
//	assert := assert.New(t)
//	for _, target := range targets {
//		err := target.AddToDatabase()
//		assert.Nil(err)
//		err = target.UpdateLatestDownloaded("newvideoid123")
//		assert.Nil(err)
//		updatedTarget, err := target.GetFromDatabase()
//		assert.Nil(err)
//		assert.Equal("newvideoid123", updatedTarget.LatestDownloaded, "updatedTarget.LatestDownloaded should be newvideoid123")
//		err = updatedTarget.Delete()
//		assert.Nil(err)
//	}
//}
//
//func TestUpdateCheckingInterval(t *testing.T) {
//	assert := assert.New(t)
//	for _, target := range targets {
//		err := target.AddToDatabase()
//		assert.Nil(err)
//		_, err = target.UpdateCheckingInterval("5")
//		assert.Nil(err)
//		updatedTarget, err := target.GetFromDatabase()
//		assert.Equal("5", updatedTarget.CheckingInterval, "updatedTarget.CheckingInterval should be 5")
//		err = updatedTarget.Delete()
//		assert.Nil(err)
//	}
//}
//
//func TestGetCheckingInterval(t *testing.T) {
//	assert := assert.New(t)
//	for _, target := range targets {
//		if target.Name != "Heroic Motivational Anthems" {
//			err := target.AddToDatabase()
//			assert.Nil(err)
//			_, err = target.UpdateCheckingInterval("5")
//			assert.Nil(err)
//			checkingInterval, err := GetCheckingInterval("channels")
//			assert.Nil(err)
//			assert.Equal(5, checkingInterval, "checkingInterval should be 5 for"+target.Name)
//			err = target.Delete()
//			assert.Nil(err)
//		}
//	}
//}
//
//func TestDownloadTarget_UpdateLastChecked(t *testing.T) {
//	currTime := time.Now().Format("01-02-2006 15:04:05")
//	assert := assert.New(t)
//	for _, target := range targets {
//		if target.Name != "Heroic Motivational Anthems" {
//			err := target.AddToDatabase()
//			assert.Nil(err)
//			err = target.UpdateLastChecked()
//			assert.Nil(err)
//			updatedTarget, err := target.GetFromDatabase()
//			assert.Nil(err)
//			assert.Equal(currTime, updatedTarget.LastChecked, "updatedTarget.LastChecked should be "+currTime)
//			err = target.Delete()
//			assert.Nil(err)
//		}
//	}
//}
