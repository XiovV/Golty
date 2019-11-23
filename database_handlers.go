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


func UpdateCheckingInterval(interval string) (Response, error) {
	log.Info("updating checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_OPENING_DATABASE", Message: "There was an error opening channels.json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
	}

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

func GetCheckingInterval() (int, error) {
	log.Info("getting checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return 0, fmt.Errorf("GetCheckingInterval: %s", err)
	}

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

func getItemFromDatabase(dbName, targetURL string) (DownloadTarget, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
	if err != nil {
		return DownloadTarget{}, fmt.Errorf("GetFromDatabase: %s", err)
	}
	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.URL == targetURL {
			return item, nil
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

func removeItem(targetURL, dbName string) error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
	if err != nil {
		return fmt.Errorf("Delete: %s", err)
	}

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.URL == targetURL {
			db = RemoveAtIndex(db, i)
			log.Info("successfully removed channel from channels.json")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+dbName)
}


func appendToDownloadHistory(target DownloadTarget, videoId, databaseName string) error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + databaseName)
	if err != nil {
		return fmt.Errorf("UpdateDownloadHistory: %s", err)
	}

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateDownloadHistory: %s", err)
	}

	for i, item := range db {
		if item.URL == target.URL {
			db[i].DownloadHistory = append(target.DownloadHistory, videoId)
			log.Info(db)
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+databaseName)
}

func addTargetToDatabase(target DownloadTarget, databaseName string) error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + databaseName)
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}
	json.Unmarshal(byteValue, &db)

	log.Info("adding channel to DB")
	db = append(db, target)
	return writeDb(db, CONFIG_ROOT+databaseName)
}

func updateLastCheckedDateAndTime(target DownloadTarget, databaseName string) error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + databaseName)
	if err != nil {
		return fmt.Errorf("UpdateLastChecked: %s", err)
	}

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateLastChecked: %s", err)
	}

	for i, item := range db {
		log.Info("LOOPING THROUGH: ", item)
		if item.URL == target.URL {
			dt := time.Now()
			db[i].LastChecked = dt.Format("01-02-2006 15:04:05")
			log.Info("last checked date and time updated successfully")
			break
		}
	}

	return writeDb(db, CONFIG_ROOT+databaseName)
}

func updateLatestDownloadedVideoId(target DownloadTarget, videoId, databaseName string) error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + databaseName)
	if err != nil {
		return fmt.Errorf("UpdateLatestDownloaded: %s", err)
	}

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return fmt.Errorf("UpdateLatestDownloaded: %s", err)
	}

	for i, item := range db {
		if item.URL == target.URL {
			db[i].LatestDownloaded = videoId
			log.Info("latest downloaded video id updated successfully")
			break
		}
	}
	return writeDb(db, CONFIG_ROOT+databaseName)
}

func checkIfTargetExists(target DownloadTarget, databaseName string) (bool, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + databaseName)
	if err != nil {
		return false, fmt.Errorf("DoesExist: %s", err)
	}

	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.URL == target.URL {
			fmt.Println(item.URL, target.URL)
			return true, nil
		}
	}

	return false, nil
}