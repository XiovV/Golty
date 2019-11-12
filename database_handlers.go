package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func openJSONDatabase(dbName string) []byte {
	jsonFile, err := os.Open(dbName)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func GetUploadsIDFromDatabase(channelName string) (string, bool) {
	channelURL := USER_URL + channelName

	byteValue := openJSONDatabase("uploadid.json")

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
	byteValue := openJSONDatabase("uploadid.json")

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
	byteValue := openJSONDatabase("uploadid.json")

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

func UpdateLatestDownloaded(channelName, videoID string) error {
	log.Info("Updating latest downloaded video id")

	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if strings.Contains(item.ChannelURL, channelName) {
			db[i].LatestDownloaded = videoID
			log.Info("Latest downloaded video id updated successfully")
			break
		}
	}

	return writeDb(db, "channels.json")
}

func AddChannelToDatabase(channelURL, downloadMode, channelName string) {
	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	log.Info("Adding channel to DB")
	db = append(db, Channel{ChannelURL: channelURL, LatestDownloaded: "", DownloadMode: downloadMode, ChannelName: channelName})
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

func writeDb(db []Channel, dbName string) error {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("There was an error writing to database: %v", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile(dbName, file, 0644)

	return nil
}

func DeleteChannel(channelURL string) {
	byteValue := openJSONDatabase("channels.json")

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
	byteValue := openJSONDatabase("channels.json")

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

func InsertFailedDownload(videoId string) error {
	byteValue := openJSONDatabase("failed.json")

	var db []FailedVideo
	var same bool
	json.Unmarshal(byteValue, &db)

	for _, video := range db {
		fmt.Println(video.VideoID, videoId)
		if video.VideoID == videoId {
			same = true
		} else {
			same = false
		}
	}
	if same == true {
		log.Info("This videoId is already in the list")
		return fmt.Errorf("This videoId is already in the list")
	}
	log.Info("Adding video id to the list")
	db = append(db, FailedVideo{VideoID: videoId})
	return fmt.Errorf("Something went terribly wrong")
}

func writeFailedVideosDb(db []FailedVideo, dbName string) error {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error marshalling json from failed.json: ", err)
		return fmt.Errorf("There was an error marshalling json from failed.json: : %v", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	err = ioutil.WriteFile(dbName, file, 0644)
	if err != nil {
		log.Error("There was an error writing to failed.json: ", err)
		return fmt.Errorf("There was an error writing to failed.json: %v", err)
	}

	return nil
}
