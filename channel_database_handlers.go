package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
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

func GetCheckingInterval() (int, error) {
	log.Info("getting checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return 0, fmt.Errorf("GetCheckingInterval: %s", err)
	}

	var db []DownloadTarget

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


