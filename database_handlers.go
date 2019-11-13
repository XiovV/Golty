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

func (c Channel) UpdateLatestDownloaded(videoID string) error {
	log.Info("Updating latest downloaded video id")

	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db[i].LatestDownloaded = videoID
			log.Info("Latest downloaded video id updated successfully")
			break
		}
	}

	return writeDb(db, "channels.json")
}

func (c Channel) GetFromDatabase() Channel {
	byteValue := openJSONDatabase("channels.json")
	var db []Channel
	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.ChannelURL == c.ChannelURL {
			return item
		}
	}

	return Channel{}
}

func (c Channel) AddToDatabase() {
	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	log.Info("Adding channel to DB")
	db = append(db, c)
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

func (c Channel) Delete() {
	log.Info("Removing channel from database")
	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db = RemoveAtIndex(db, i)
			log.Info("Successfully removed channel from channels.json")
		}
	}

	writeDb(db, "channels.json")
}

func (c Channel) DoesExist() bool {
	byteValue := openJSONDatabase("channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for _, channel := range db {
		if channel.ChannelURL == c.ChannelURL {
			return true
		}
	}

	return false
}
