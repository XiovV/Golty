package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (v DownloadVideoPayload) AddToDatabase() error {
	b, err := ioutil.ReadFile(filepath.Join(CONFIG_ROOT, "videos.json"))
	if err != nil {
		return fmt.Errorf("v.AddToDatabase: %s", err)
	}

	var videos []DownloadVideoPayload

	json.Unmarshal(b, &videos)

	log.Info("adding video to DB")
	videos = append(videos, v)
	err = writeToVideosDb(videos, CONFIG_ROOT+"videos.json")
	if err != nil {
		return fmt.Errorf("v.AddToDatabase: %s", err)
	}
	return nil
}
func GetVideos() ([]DownloadVideoPayload, error) {
	log.Info("getting all videos from videos.json")
	jsonFile, err := os.Open(CONFIG_ROOT + "videos.json")
	if err != nil {
		log.Error("From GetChannels()", err)
		return []DownloadVideoPayload{}, fmt.Errorf("From GetVideos(): %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var videos []DownloadVideoPayload

	err = json.Unmarshal(byteValue, &videos)
	if err != nil {
		log.Error("From GetChannels()", err)
		return []DownloadVideoPayload{}, fmt.Errorf("From GetChannels(): %v", err)
	}
	log.Info("successfully read all channels")
	return reverseVideos(videos), nil
}

func writeToVideosDb(db []DownloadVideoPayload, dbName string) error {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("writeToVideosDb: %s", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	err = ioutil.WriteFile(dbName, file, 0644)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("writeToVideosDb: %s", err)
	}

	return nil
}
