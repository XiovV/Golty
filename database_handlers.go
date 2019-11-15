package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

func openJSONDatabase(dbName string) ([]byte, error) {
	jsonFile, err := os.Open(dbName)
	if err != nil {
		log.Errorf("openJSONDatabase: %s", err)
		return nil, fmt.Errorf("openJSONDatabase: %s", err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("openJSONDatabase: %s", err)
		return nil, fmt.Errorf("openJSONDatabase: %s", err)
	}

	return byteValue, nil
}

func (c *Channel) UpdateDownloadHistory(videoID string) error {
	log.Info("updating download history")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return fmt.Errorf("UpdateDownloadHistory: %s", err)
	}

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateDownloadHistory: %s", err)
	}

	for i, channel := range db {
		if channel.ChannelURL == c.ChannelURL {
			db[i].DownloadHistory = append(c.DownloadHistory, videoID)
			// log.Info(channel.ChannelURL, channel.LatestDownloaded, channel.DownloadMode, channel.Name, channel.PreferredExtensionForAudio, channel.PreferredExtensionForVideo, c.DownloadHistory)
			log.Info(db)
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+"channels.json")
}

func (c Channel) UpdateLatestDownloaded(videoID string) error {
	log.Info("updating latest downloaded video id")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return fmt.Errorf("UpdateLatestDownloaded: %s", err)
	}

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateLatestDownloaded: %s", err)
	}

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db[i].LatestDownloaded = videoID
			log.Info("Latest downloaded video id updated successfully")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+"channels.json")
}

func (c Channel) GetFromDatabase() (Channel, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return Channel{}, fmt.Errorf("GetFromDatabase: %s", err)
	}
	var db []Channel
	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.ChannelURL == c.ChannelURL {
			return item, nil
		}
	}

	return Channel{}, fmt.Errorf("Couldn't find channel in the database: %s", c.ChannelURL)
}

func (c Channel) AddToDatabase() error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}

	var db []Channel

	json.Unmarshal(byteValue, &db)

	log.Info("Adding channel to DB")
	db = append(db, c)
	err = writeDb(db, CONFIG_ROOT+"channels.json")
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}
	return nil
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
		return fmt.Errorf("writeDb: %s", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	err = ioutil.WriteFile(dbName, file, 0644)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("writeDb: %s", err)
	}

	return nil
}

func (c Channel) Delete() error {
	log.Info("Removing channel from database")
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return fmt.Errorf("Delete: %s", err)
	}
	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			db = RemoveAtIndex(db, i)
			log.Info("successfully removed channel from channels.json")
		}
	}

	writeDb(db, CONFIG_ROOT+"channels.json")
	return nil
}

func (c Channel) DoesExist() (bool, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return false, fmt.Errorf("AddToDatabase: %s", err)
	}
	var db []Channel

	json.Unmarshal(byteValue, &db)

	for _, channel := range db {
		if channel.ChannelURL == c.ChannelURL {
			return true, nil
		}
	}

	return false, nil
}
