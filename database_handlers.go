package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

func (c Channel) UpdateLatestDownloaded(videoID string) error {
	log.Info("Updating latest downloaded video id")

	byteValue := openJSONDatabase(CONFIG_ROOT + "channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db[i].LatestDownloaded = videoID
			log.Info("Latest downloaded video id updated successfully")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+"channels.json")
}

func (c Channel) GetFromDatabase() Channel {
	byteValue := openJSONDatabase(CONFIG_ROOT + "channels.json")
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
	byteValue := openJSONDatabase(CONFIG_ROOT + "channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	log.Info("Adding channel to DB")
	db = append(db, c)
	writeDb(db, CONFIG_ROOT+"channels.json")
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
	byteValue := openJSONDatabase(CONFIG_ROOT + "channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db = RemoveAtIndex(db, i)
			log.Info("successfully removed channel from channels.json")
		}
	}

	writeDb(db, CONFIG_ROOT+"channels.json")
}

func (c Channel) DoesExist() bool {
	byteValue := openJSONDatabase(CONFIG_ROOT + "channels.json")

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for _, channel := range db {
		if channel.ChannelURL == c.ChannelURL {
			return true
		}
	}

	return false
}
