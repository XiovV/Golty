package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func UpdateFailedVideos(videoId string) {
	jsonFile, err := os.Open("failed.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []FailedVideos

	json.Unmarshal(byteValue, &db)

	db = append(db, FailedVideos{VideoID: videoId})
	writeFailedVideos(db, "failed.json")
}

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
			return item.UploadsID, true
			break
		}
	}

	return "", false
}

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

func UpdateUploadsIDDatabase(channelURL string) bool {
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if strings.Contains(item.ChannelURL, channelName) {
			db[i].LatestDownloaded = videoID
			break
		}
	}

	writeDb(db, "channels.json")
}

func UpdateChannelsDatabase(channelURL string) bool {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

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
		db = append(db, Channel{ChannelURL: channelURL})
		writeDb(db, "channels.json")
	}
	return false
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

func writeFailedVideos(db []FailedVideos, dbName string) {
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
		log.Println(err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile(dbName, file, 0644)
}

func DeleteChannel(channelURL string) {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == channelURL {
			db = RemoveAtIndex(db, i)
		}
	}

	writeDb(db, "channels.json")
}
