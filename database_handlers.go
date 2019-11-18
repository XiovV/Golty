package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"time"
)

func UpdateCheckingInterval(interval string) (Response, error) {
	log.Info("updating checking interval")

	byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_OPENING_DATABASE", Message: "There was an error opening channels.json: " + err.Error()}, fmt.Errorf("UpdateCheckingInterval: %s", err)
	}

	var db []DownloadTarget

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



//func DoesExist(target interface{}) (bool, error) {
//	switch target := target.(type) {
//	case Channel:
//		byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
//		if err != nil {
//			return false, fmt.Errorf("DoesExist: %s", err)
//		}
//		var db []Channel
//
//		json.Unmarshal(byteValue, &db)
//
//		for _, channel := range db {
//			if channel.URL == target.URL {
//				fmt.Println(channel.URL, target.URL)
//				return true, nil
//			}
//		}
//		return false, nil
//	case Playlist:
//		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
//		if err != nil {
//			return false, fmt.Errorf("DoesExist: %s", err)
//		}
//		var playlists []Playlist
//
//		json.Unmarshal(byteValue, &playlists)
//
//		for _, playlist := range playlists {
//			if playlist.URL == target.URL {
//				return true, nil
//			}
//		}
//		return false, nil
//	}
//	return true, nil
//}

//func UpdateLastChecked(target interface{}) error {
//	switch target := target.(type) {
//	case Channel:
//		log.Info("updating last checked date and time for: ", target.Name)
//
//		byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
//		if err != nil {
//			return fmt.Errorf("UpdateLastChecked: %s", err)
//		}
//
//		var db []Channel
//
//		err = json.Unmarshal(byteValue, &db)
//		if err != nil {
//			return fmt.Errorf("UpdateLastChecked: %s", err)
//		}
//
//		for i, item := range db {
//			if item.URL == target.URL {
//				dt := time.Now()
//				db[i].LastChecked = dt.Format("01-02-2006 15:04:05")
//				log.Info("last checked date and time updated successfully")
//				break
//			}
//		}
//
//	case Playlist:
//		log.Info("updating last checked date and time for: ", target.Name)
//
//		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
//		if err != nil {
//			return fmt.Errorf("p.UpdateLastChecked: %s", err)
//		}
//
//		var playlists []Playlist
//
//		err = json.Unmarshal(byteValue, &playlists)
//		if err != nil {
//			return fmt.Errorf("p.UpdateLastChecked: %s", err)
//		}
//
//		for i, playlist := range playlists {
//			if playlist.URL == target.URL {
//				dt := time.Now()
//				playlists[i].LastChecked = dt.Format("01-02-2006 15:04:05")
//				log.Info("last checked date and time updated successfully")
//				break
//			}
//		}
//
//		return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
//	}
//
//	return nil
//}

func UpdateLastChecked(target DownloadTarget) error {
	log.Info("UPDATING LAST CHECKED FOR: ", target.URL)
	var db []DownloadTarget
	var dbName string
	log.Info("updating last checked date and time for: ", target.Name)
	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}

	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
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

	return writeDb(db, CONFIG_ROOT+dbName)
}

func DoesExist(target DownloadTarget) (bool, error) {
	var db []DownloadTarget
	var dbName string

	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}
	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
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

func UpdateLatestDownloaded(target DownloadTarget, videoId string) error {
	log.Info("updating latest downloaded video id")
	var db []DownloadTarget
	var dbName string
	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}

	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
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
	return writeDb(db, CONFIG_ROOT+dbName)
}

func AddToDatabase(target DownloadTarget) error {
	var db []DownloadTarget
	var dbName string
	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}
	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}
	json.Unmarshal(byteValue, &db)

	log.Info("adding channel to DB")
	db = append(db, target)
	err = writeDb(db, CONFIG_ROOT+dbName)
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}

	return nil
}

func UpdateDownloadHistory(target DownloadTarget, videoId string) error {
	log.Info("updating download history")
	var db []DownloadTarget
	var dbName string
	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}

	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
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

	return writeDb(db, CONFIG_ROOT+dbName)
}

func Delete(target DownloadTarget) error {
	var db []DownloadTarget
	var dbName string
	if target.Type == "Channel" {
		log.Info("Removing channel from database")
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		log.Info("Removing playlist from database")
		dbName = "playlists.json"
	}

	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
	if err != nil {
		return fmt.Errorf("Delete: %s", err)
	}

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if item.URL == target.URL {
			db = RemoveAtIndex(db, i)
			log.Info("successfully removed channel from channels.json")
		}
	}

	return writeDb(db, CONFIG_ROOT+dbName)
}

func GetFromDatabase(target DownloadTarget) (DownloadTarget, error) {
	var db []DownloadTarget
	var dbName string
	if target.Type == "Channel" {
		dbName = "channels.json"
	} else if target.Type == "Playlist" {
		dbName = "playlists.json"
	}
	byteValue, err := openJSONDatabase(CONFIG_ROOT + dbName)
	if err != nil {
		return DownloadTarget{}, fmt.Errorf("GetFromDatabase: %s", err)
	}
	json.Unmarshal(byteValue, &db)

	for _, item := range db {
		if item.URL == target.URL {
			return item, nil
		}
	}

	return DownloadTarget{}, fmt.Errorf("Couldn't find channel/playlist in the database: %s", target.URL)
}