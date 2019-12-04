package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var db []DownloadTarget
var databaseName string

func (target DownloadTarget) UpdateCheckingInterval(interval string) (Response, error) {
	databaseName = setDatabaseName(target.Type)

	log.Info("updating checking interval")

	db = openDatabaseAndUnmarshalJSON(databaseName)

	if len(db) > 0 {
		db[0].CheckingInterval = interval
		err := writeDb(db, CONFIG_ROOT+databaseName)
		if err != nil {
			return Response{Type: "Error", Key: "ERROR_WRITING_TO_DATABASE", Message: "There was an error writing to channels.json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
		}
		return Response{Type: "Success", Key: "UPDATE_CHECKING_INTERVAL_SUCCESS", Message: "Successfully updated the checking interval"}, nil
	}
	return Response{Type: "Error", Key: "DATABASE_EMPTY", Message: "There has to be at least one channel in the database before updating the checking interval."}, nil
}

func (target DownloadTarget) UpdateLastChecked() error {
	log.Info("UPDATING LAST CHECKED FOR: ", target.URL)
	log.Info("updating last checked date and time for: ", target.Name)
	databaseName = setDatabaseName(target.Type)
	return updateLastCheckedDateAndTime(target, databaseName)
}

func (target DownloadTarget) DoesExist() (bool, error) {
	databaseName = setDatabaseName(target.Type)
	return checkIfTargetExists(target, databaseName)
}

func (target DownloadTarget) UpdateLatestDownloaded(videoId string) error {
	log.Info("updating latest downloaded video id")
	databaseName = setDatabaseName(target.Type)
	return updateLatestDownloadedVideoId(target, videoId, databaseName)
}

func (target DownloadTarget) AddToDatabase() error {
	databaseName = setDatabaseName(target.Type)
	return addTargetToDatabase(target, databaseName)
}

func (target DownloadTarget) UpdateDownloadHistory(videoId string) error {
	log.Info("updating download history")
	databaseName = setDatabaseName(target.Type)
	return appendToDownloadHistory(target, videoId, databaseName)
}


func (target DownloadTarget) Delete() error {
	databaseName = setDatabaseName(target.Type)
	return removeItem(target.URL, databaseName)
}

func (target DownloadTarget) GetFromDatabase() (DownloadTarget, error) {
	databaseName = setDatabaseName(target.Type)
	return getItemFromDatabase(databaseName, target.URL)
}

func GetCheckingInterval(target string) (int, error) {
	log.Info("getting checking interval")

	if target == "channels" {
		db = openDatabaseAndUnmarshalJSON("channels.json")
		if len(db) > 0 && db[0].CheckingInterval != "" {
			checkingInterval, err := strconv.Atoi(db[0].CheckingInterval)
			if err != nil {
				return 0, fmt.Errorf("GetCheckingInterval: %s", err)
			}
			log.Info("got checking interval successfully")
			return checkingInterval, nil
		}
	} else if target == "playlists" {
		db = openDatabaseAndUnmarshalJSON("playlists.json")
		if len(db) > 0 && db[0].CheckingInterval != "" {
			checkingInterval, err := strconv.Atoi(db[0].CheckingInterval)
			if err != nil {
				return 0, fmt.Errorf("GetCheckingInterval: %s", err)
			}
			log.Info("got checking interval successfully")
			return checkingInterval, nil
		}
	}

	log.Info("checking interval not yet specified")
	return 0, nil
}

func writeDb(db []DownloadTarget, dbName string) error {
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

func getItemFromDatabase(databaseName, targetURL string) (DownloadTarget, error) {
	db = openDatabaseAndUnmarshalJSON(databaseName)

	for i := range db {
		if db[i].URL == targetURL {
			return db[i], nil
		}
	}
	return DownloadTarget{}, fmt.Errorf("Couldn't find target")
}


func setDatabaseName(targetType string) string {
	if targetType == "Channel" {
		log.Info("Removing channel from database")
		return "channels.json"
	} else if targetType == "Playlist" {
		log.Info("Removing playlist from database")
		return "playlists.json"
	}
	return ""
}

func removeItem(targetURL, databaseName string) error {
	db = openDatabaseAndUnmarshalJSON(databaseName)

	for i := range db {
		if db[i].URL == targetURL {
			db = RemoveAtIndex(db, i)
			log.Info("successfully removed channel from channels.json")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+databaseName)
}


func appendToDownloadHistory(target DownloadTarget, videoId, databaseName string) error {
	db = openDatabaseAndUnmarshalJSON(databaseName)

	for i:= range db {
		if db[i].URL == target.URL {
			db[i].DownloadHistory = append(target.DownloadHistory, videoId)
			log.Info(db)
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+databaseName)
}

func addTargetToDatabase(target DownloadTarget, databaseName string) error {
	openDatabaseAndUnmarshalJSON(databaseName)

	log.Info("adding channel to DB")
	db = append(db, target)
	return writeDb(db, CONFIG_ROOT+databaseName)
}

func updateLastCheckedDateAndTime(target DownloadTarget, databaseName string) error {
	openDatabaseAndUnmarshalJSON(databaseName)

	for i:= range db {
		if db[i].URL == target.URL {
			dt := time.Now()
			db[i].LastChecked = dt.Format("01-02-2006 15:04:05")
			log.Info("last checked date and time updated successfully")
			break
		}
	}
	return writeDb(db, CONFIG_ROOT+databaseName)
}

func updateLatestDownloadedVideoId(target DownloadTarget, videoId, databaseName string) error {
	db = openDatabaseAndUnmarshalJSON(databaseName)

	for i := range db {
		if db[i].URL == target.URL {
			db[i].LatestDownloaded = videoId
			log.Info("latest downloaded video id updated successfully")
			break
		}
	}
	return writeDb(db, CONFIG_ROOT+databaseName)
}

func checkIfTargetExists(target DownloadTarget, databaseName string) (bool, error) {
	db = openDatabaseAndUnmarshalJSON(databaseName)

	for i := range db {
		if db[i].URL == target.URL {
			fmt.Println(db[i].URL, target.URL)
			return true, nil
		}
	}

	return false, nil
}

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

func openDatabaseAndUnmarshalJSON(databaseName string) []DownloadTarget {
	byteValue, _ := openJSONDatabase(CONFIG_ROOT + databaseName)
	json.Unmarshal(byteValue, &db)
	return db
}