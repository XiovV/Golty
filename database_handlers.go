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

func UpdateLastChecked(target interface{}) error {
	switch target := target.(type) {
	case Channel:
		log.Info("updating last checked date and time for: ", target.Name)

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
			if item.ChannelURL == target.ChannelURL {
				dt := time.Now()
				db[i].LastChecked = dt.Format("01-02-2006 15:04:05")
				log.Info("last checked date and time updated successfully")
				break
			}
		}

	case Playlist:
		log.Info("updating last checked date and time for: ", target.Name)

		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return fmt.Errorf("p.UpdateLastChecked: %s", err)
		}

		var playlists []Playlist

		err = json.Unmarshal(byteValue, &playlists)
		if err != nil {
			return fmt.Errorf("p.UpdateLastChecked: %s", err)
		}

		for i, playlist := range playlists {
			if playlist.PlaylistURL == target.PlaylistURL {
				dt := time.Now()
				playlists[i].LastChecked = dt.Format("01-02-2006 15:04:05")
				log.Info("last checked date and time updated successfully")
				break
			}
		}

		return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
	}

	return nil
}

func AddToDatabase(target interface{}) error {
	switch target := target.(type) {
	case Channel:
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
		if err != nil {
			return fmt.Errorf("AddToDatabase: %s", err)
		}

		var db []Channel

		json.Unmarshal(byteValue, &db)

		log.Info("adding channel to DB")
		db = append(db, target)
		err = writeDb(db, CONFIG_ROOT+"channels.json")
		if err != nil {
			return fmt.Errorf("AddToDatabase: %s", err)
		}
		return nil

	case Playlist:
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return fmt.Errorf("AddToDatabase: %s", err)
		}

		var playlists []Playlist

		json.Unmarshal(byteValue, &playlists)

		log.Info("adding channel to DB")
		playlists = append(playlists, target)
		err = writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
		if err != nil {
			return fmt.Errorf("AddToDatabase: %s", err)
		}
		return nil
	}
	return nil
}

func DoesExist(target interface{}) (bool, error) {
	switch target := target.(type) {
	case Channel:
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
		if err != nil {
			return false, fmt.Errorf("DoesExist: %s", err)
		}
		var db []Channel

		json.Unmarshal(byteValue, &db)

		for _, channel := range db {
			if channel.ChannelURL == target.ChannelURL {
				fmt.Println(channel.ChannelURL, target.ChannelURL)
				return true, nil
			}
		}
		return false, nil
	case Playlist:
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return false, fmt.Errorf("DoesExist: %s", err)
		}
		var playlists []Playlist

		json.Unmarshal(byteValue, &playlists)

		for _, playlist := range playlists {
			if playlist.PlaylistURL == target.PlaylistURL {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}

func UpdateLatestDownloaded(target interface{}, videoId string) error {
	log.Info("updating latest downloaded video id")

	switch target := target.(type) {
	case Channel:
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
			if item.ChannelURL == target.ChannelURL {
				db[i].LatestDownloaded = videoId
				log.Info("latest downloaded video id updated successfully")
				break
			}
		}
		return writeDb(db, CONFIG_ROOT+"channels.json")

	case Playlist:
		log.Info("updating latest downloaded video id")

		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return fmt.Errorf("p.UpdateLatestDownloaded: %s", err)
		}

		var playlists []Playlist

		err = json.Unmarshal(byteValue, &playlists)
		if err != nil {
			return fmt.Errorf("p.UpdateLatestDownloaded: %s", err)
		}

		for i, playlist := range playlists {
			if playlist.PlaylistURL == target.PlaylistURL {
				playlists[i].LatestDownloaded = videoId
				log.Info("latest downloaded video id updated successfully")
				break
			}
		}

		return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
	}

	return nil
}

func UpdateDownloadHistory(target interface{}, videoId string) error {
	log.Info("updating download history")
	switch target := target.(type) {
	case Channel:
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
			if channel.ChannelURL == target.ChannelURL {
				db[i].DownloadHistory = append(target.DownloadHistory, videoId)
				log.Info(db)
				break
			}
		}

		return writeDb(db, CONFIG_ROOT+"channels.json")

	case Playlist:
		log.Info("updating download history")

		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return fmt.Errorf("p.UpdateDownloadHistory: %s", err)
		}

		var playlists []Playlist

		err = json.Unmarshal(byteValue, &playlists)
		if err != nil {
			return fmt.Errorf("p.UpdateDownloadHistory: %s", err)
		}

		for i, playlist := range playlists {
			if playlist.PlaylistURL == target.PlaylistURL {
				playlists[i].DownloadHistory = append(target.DownloadHistory, videoId)
				log.Info(playlists)
				break
			}
		}

		return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
	}

	return nil
}

func Delete(target interface{}) error {
	switch target := target.(type) {
	case Channel:
		log.Info("Removing channel from database")
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "channels.json")
		if err != nil {
			return fmt.Errorf("Delete: %s", err)
		}
		var db []Channel

		json.Unmarshal(byteValue, &db)

		for i, item := range db {
			if item.ChannelURL == target.ChannelURL {
				db = RemoveAtIndexChannel(db, i)
				log.Info("successfully removed channel from channels.json")
			}
		}

		writeDb(db, CONFIG_ROOT+"channels.json")
		return nil

	case Playlist:
		log.Info("Removing channel from database")
		byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
		if err != nil {
			return fmt.Errorf("p.Delete: %s", err)
		}
		var playlists []Playlist

		json.Unmarshal(byteValue, &playlists)

		for i, playlist := range playlists {
			if playlist.PlaylistURL == target.PlaylistURL {
				playlists = RemoveAtIndexPlaylist(playlists, i)
				log.Info("successfully removed playlist from playlists.json")
			}
		}

		writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
		return nil
	}
	return nil
}
