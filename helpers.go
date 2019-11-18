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
	var dbName string
	if target == "channels" {
		dbName = "channels.json"
	} else if target == "playlists" {
		dbName = "playlists.json"
	}

	log.Info("getting all channels from channels.json")
	jsonFile, err := os.Open(CONFIG_ROOT + dbName)
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
	return db, nil
}

func CheckAll(target string) (Response, error) {
	var foundFor []string
	var preferredExtension string
	var allInInDb []DownloadTarget
	var err error
	var targetType string
	if target == "channels" {
		log.Info("checking for all channels")
		allInInDb, err = GetAll("channels")
		targetType = "Channel"
	} else if target == "playlists" {
		log.Info("checking for all playlists")
		allInInDb, err = GetAll("playlists")
		targetType = "Playlist"
	}
	for _, item := range allInInDb {
		target := DownloadTarget{URL: item.URL, Type: targetType}
		target, err = GetFromDatabase(target)
		if err != nil {
			return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the channel from database" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
		}

		if item.URL == target.URL {
			videoId, err := GetLatestVideo(target)
			if err != nil {
				log.Error("There was an error getting latest video: ", err)
				return Response{Type: "Error", Key: "GETTING_LATEST_VIDEO_ERROR", Message: "There was an error getting the latestvideo" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
			}

			UpdateLastChecked(item)
			if item.LatestDownloaded == videoId {
				log.Info("no new videos found for: ", item.URL)
			} else {
				log.Info("new video detected for: ", item.URL)
				foundFor = append(foundFor, item.URL)
				if target.DownloadMode == "Audio Only" {
					preferredExtension = target.PreferredExtensionForAudio
				} else if target.DownloadMode == "Video And Audio" {
					preferredExtension = target.PreferredExtensionForVideo
				}
				go Download(target, "best", preferredExtension, false)
				UpdateLatestDownloaded(target, videoId)
			}
		}
	}
	if len(foundFor) == 0 {
		return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos found."}, nil
	}
	return Response{Type: "Success", Key: "NEW_VIDEOS", Message: strings.Join(foundFor, ",")}, nil
}

func (target DownloadTarget) CheckNow() (Response, error) {
	var allInDb []DownloadTarget
	var err error
	log.Info("checking for new videos")
	if target.Type == "Channel" {
		allInDb, err = GetAll("channels")
		if err != nil {
			return Response{}, fmt.Errorf("From p.CheckNow(): %v", err)
		}
	} else if target.Type == "Playlist" {
		allInDb, err = GetAll("playlists")
		if err != nil {
			return Response{}, fmt.Errorf("From p.CheckNow(): %v", err)
		}
	}
	var preferredExtension string

	targets, err := GetFromDatabase(target)
	if err != nil {
		return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the playlist from database" + err.Error()}, fmt.Errorf("CheckNow: %s", err)
	}
	targetURL := target.URL

	targetMetadata, err := GetMetadata(targets)
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting playlist metadata: " + err.Error()}, nil
	}

	err = UpdateLastChecked(target)
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_UPDATING_LAST_CHECKED", Message: "There was an error updating latest checked date and time: " + err.Error()}, nil
	}

	for _, target := range allInDb {
		if target.URL == targetURL {
			if target.LatestDownloaded == targetMetadata.ID {
				log.Info("no new videos found for: ", targetURL)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected for " + target.Name}, nil
			} else {
				log.Info("new video detected for: ", targetURL)
				if target.DownloadMode == "Audio Only" {
					preferredExtension = target.PreferredExtensionForAudio
				} else if target.DownloadMode == "Video And Audio" {
					preferredExtension = target.PreferredExtensionForVideo
				}
				err := Download(target, "best", preferredExtension, false)
				if err != nil {
					log.Error(err)
					return Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: err.Error()}, nil
				}
				UpdateLatestDownloaded(target, targetMetadata.ID)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected for " + target.Name + " and downloaded"}, nil
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}, nil
}

func RemoveAtIndex(s []DownloadTarget, index int) []DownloadTarget {
	return append(s[:index], s[index+1:]...)
}

func GetChannelName(channelURL string) string {
	return strings.Split(channelURL, "/")[4]
}

func ReturnResponse(w http.ResponseWriter, res Response) {
	log.Info("returning response: ", res)
	json.NewEncoder(w).Encode(res)
}