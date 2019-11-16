package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

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

func (c Channel) UpdateDownloadHistory(videoID string) error {
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
			log.Info(db)
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+"channels.json")
}

func (c Channel) UpdateLastChecked() error {
	log.Info("updating last checked date and time for: ", c.Name)

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return fmt.Errorf("UpdateLastChecked: %s", err)
	}

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateLastChecked: %s", err)
	}

	for i, item := range db {
		if item.ChannelURL == c.ChannelURL {
			dt := time.Now()
			db[i].LastChecked = dt.Format("01-02-2006 15:04:05")
			log.Info("last checked date and time updated successfully")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+"channels.json")
}

func UpdateCheckingInterval(interval string) (Response, error) {
	log.Info("updating checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_OPENING_DATABASE", Message: "There was an error opening channels.json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
	}

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_UNMRASHALLING_JSON", Message: "There was an error unmarshalling json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
	}

	if len(db) > 0 {
		db[0].CheckingInterval = interval
		err = writeDb(db, CONFIG_ROOT+"channels.json")
		if err != nil {
			return Response{Type: "Error", Key: "ERROR_WRITING_TO_DATABASE", Message: "There was an error writing to channels.json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
		}
		return Response{Type: "Success", Key: "UPDATE_CHECKING_INTERVAL_SUCCESS", Message: "Successfully updated the checking interval"}, nil
	}
	return Response{Type: "Error", Key: "DATABASE_EMPTY", Message: "There has to be at least one channel in the database before updating the checking interval."}, nil
}

func GetCheckingInterval() (int, error) {
	log.Info("getting checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return 0, fmt.Errorf("GetCheckingInterval: %s", err)
	}

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return 0, fmt.Errorf("GetCheckingInterval: %s", err)
	}
	if len(db) > 0 && db[0].CheckingInterval != "" {
		checkingInterval, err := strconv.Atoi(db[0].CheckingInterval)
		if err != nil {
			return 0, fmt.Errorf("GetCheckingInterval: %s", err)
		}
		log.Info("got checking interval successfully")
		return checkingInterval, nil
	}
	log.Info("checking interval not yet specified")
	return 0, nil
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
			log.Info("latest downloaded video id updated successfully")
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

	log.Info("adding channel to DB")
	db = append(db, c)
	err = writeDb(db, CONFIG_ROOT+"channels.json")
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}
	return nil
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
			db = RemoveAtIndexChannel(db, i)
			log.Info("successfully removed channel from channels.json")
		}
	}

	writeDb(db, CONFIG_ROOT+"channels.json")
	return nil
}

func (c Channel) DoesExist() (bool, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return false, fmt.Errorf("DoesExist: %s", err)
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
