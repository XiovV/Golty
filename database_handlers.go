package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetUploadsIDFromDatabase(channelName string) (string, bool) {
	channelURL := USER_URL + channelName

	jsonFile, err := os.Open("uploadid.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []UploadID

	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.ChannelURL == channelURL {
			log.Info("Found uploads id for this user")
			return item.UploadsID, true
		}
	}

	return "", false
}

// UpdateUploadsID stores/updates UploadsID inside uploadid.json for a channel
func UpdateUploadsID(channelName, uploadsId string) {
	jsonFile, err := os.Open("uploadid.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []UploadID

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if strings.Contains(item.ChannelURL, channelName) {
			db[i].UploadsID = uploadsId
			break
		}
	}

	writeUploadsDb(db, "uploadid.json")
}

// InitUploadsID puts a channel inside uploadid.json and leaves UploadsID empty.
func InitUploadsID(channelURL string) bool {
	jsonFile, err := os.Open("uploadid.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []UploadID

	json.Unmarshal(byteValue, &db)

	var exists bool

	for _, v := range db {
		if v.ChannelURL == channelURL {
			exists = true
			break
		} else if channelURL == "" {
			fmt.Println("channelURL can't be empty", channelURL)
			exists = true
			break
		} else {
			exists = false
		}
	}

	if exists == true {
		fmt.Println("channel already added", channelURL)
		return true
	} else {
		fmt.Println("adding to db: ", channelURL)
		db = append(db, UploadID{ChannelURL: channelURL})
		writeUploadsDb(db, "uploadid.json")
	}
	return false
}

func UpdateLatestDownloaded(channelName, videoID string) {
	log.Info("Updating latest downloaded video id")
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Error("There was an error reading channels.json: ", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if strings.Contains(item.ChannelURL, channelName) {
			db[i].LatestDownloaded = videoID
			log.Info("Latest downloaded video id updated successfully")
			break
		}
	}

	writeDb(db, "channels.json")
}

func AddChannelToDatabase(channelURL string) {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	log.Info("Adding channel to DB")
	db = append(db, Channel{ChannelURL: channelURL})
	writeDb(db, "channels.json")
}

func writeUploadsDb(db []UploadID, dbName string) {
	result, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile(dbName, file, 0644)
}

func writeDb(db []Channel, dbName string) {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile(dbName, file, 0644)
}

func DeleteChannel(channelURL string) {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Error("There was an error reading channels.json", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == channelURL {
			db = RemoveAtIndex(db, i)
			log.Info("Successfully removed channel from channels.json")
		}
	}

	writeDb(db, "channels.json")
}

func DoesChannelExist(channelURL string) bool {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for _, channel := range db {
		if channel.ChannelURL == channelURL {
			return true
		}
	}

	return false
}

func CheckIfDownloadFailed(videoPath string) bool {
	log.Info("Checking if the download failed")
	path := strings.Split(videoPath, "/")
	videoTitle := path[1]
	files, err := ioutil.ReadDir(path[0])
	if err != nil {
		log.Error("There was a problem reading the directory: ", err)
	}

	for _, f := range files {
		if f.Name() == videoTitle {
			return false
		}
	}

	return true
}

func InsertFailedDownload(videoId string) {
	jsonFile, err := os.Open("failed.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []FailedVideo

	json.Unmarshal(byteValue, &db)

	for _, video := range db {
		if video.VideoID == videoId {
			log.Info("This videoId is already in the list")
		} else {
			log.Info("Adding channel to DB")
			db = append(db, FailedVideo{VideoID: videoId})
			writeFailedVideosDb(db, "failed.json")
		}
	}
}

func writeFailedVideosDb(db []FailedVideo, dbName string) {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile(dbName, file, 0644)
}
