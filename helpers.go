package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetAll(target string) ([]DownloadTarget, error) {
	var db []DownloadTarget
	var databaseName string
	if target == "channels" {
		databaseName = "channels.json"
	} else if target == "playlists" {
		databaseName = "playlists.json"
	}

	log.Info("getting all channels from channels.json")
	jsonFile, err := os.Open(CONFIG_ROOT + databaseName)
	if err != nil {
		log.Error("From GetAll()", err)
		return []DownloadTarget{}, fmt.Errorf("From GetAll(): %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		log.Error("From GetAll()", err)
		return []DownloadTarget{}, fmt.Errorf("From GetAll(): %v", err)
	}
	log.Info("successfully read all channels")
	log.Info(db)
	return db, nil
}

func CheckAll(target string) (Response, error) {
	var foundFor []string
	var preferredExtension string
	var allInDb []DownloadTarget
	var err error
	var targetType string
	if target == "channels" {
		log.Info("checking for all channels")
		allInDb, err = GetAll("channels")
		targetType = "Channel"
	} else if target == "playlists" {
		log.Info("checking for all playlists")
		allInDb, err = GetAll("playlists")
		targetType = "Playlist"
	}
	for _, item := range allInDb {
		target := DownloadTarget{URL: item.URL, Type: targetType}
		target, err = target.GetFromDatabase()
		if err != nil {
			return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the channel from database" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
		}
		newVideoUploaded, videoId, err := target.CheckNow()
		if err != nil {
			return Response{Type: "Error", Key: "CHECKING_ERROR", Message: "There was an error checking for new uploads" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
		}
		if newVideoUploaded == true {
			log.Info("new video detected for: ", item.URL)
			foundFor = append(foundFor, item.URL)
			if target.DownloadMode == "Audio Only" {
				preferredExtension = target.PreferredExtensionForAudio
			} else if target.DownloadMode == "Video And Audio" {
				preferredExtension = target.PreferredExtensionForVideo
			}
			go target.Download("best", preferredExtension, false)
			target.UpdateLatestDownloaded(videoId)
		} else {
			log.Info("no new videos found for: ", item.URL)
		}
	}
	if len(foundFor) == 0 {
		return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos found."}, nil
	}
	return Response{Type: "Success", Key: "NEW_VIDEOS", Message: "New videos detected for: " + strings.Join(foundFor, ",")}, nil
}

func (target DownloadTarget) CheckNow() (bool, string, error) {
	var allInDb []DownloadTarget
	var err error
	log.Info("checking for new videos")
	if target.Type == "Channel" {
		allInDb, err = GetAll("channels")
		if err != nil {
			return false, "", fmt.Errorf("From p.CheckNow(): %v", err)
		}
	} else if target.Type == "Playlist" {
		allInDb, err = GetAll("playlists")
		if err != nil {
			return false, "", fmt.Errorf("From p.CheckNow(): %v", err)
		}
	}

	targets, err := target.GetFromDatabase()
	if err != nil {
		return false, "", fmt.Errorf("CheckNow: %s", err)
	}
	targetURL := target.URL

	targetMetadata, err := targets.GetMetadata()
	if err != nil {
		return false, "", fmt.Errorf("CheckNow: %s", err)
	}

	err = target.UpdateLastChecked()
	if err != nil {
		return false, "", fmt.Errorf("CheckNow: %s", err)
	}

	for _, target := range allInDb {
		if target.URL == targetURL {
			if target.LatestDownloaded == targetMetadata.ID {
				log.Info("no new videos found for: ", targetURL)
				return false, "", nil
			} else {
				log.Info("new video found for: ", targetURL)
				return true, targetMetadata.ID, nil
			}
		}
	}
	log.Error("CheckNow: Something went terribly wrong")
	return false, "", fmt.Errorf("CheckNow: something went wrong ")
}

func RemoveAtIndex(s []DownloadTarget, index int) []DownloadTarget {
	return append(s[:index], s[index+1:]...)
}

func ReturnResponse(w http.ResponseWriter, res Response) {
	log.Info("returning response: ", res)
	json.NewEncoder(w).Encode(res)
}

func Log(err error) error {
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}