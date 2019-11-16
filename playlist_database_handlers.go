package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

func (p Playlist) DoesExist() (bool, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
	if err != nil {
		return false, fmt.Errorf("DoesExist: %s", err)
	}
	var playlists []Playlist

	json.Unmarshal(byteValue, &playlists)

	for _, playlist := range playlists {
		if playlist.PlaylistURL == p.PlaylistURL {
			return true, nil
		}
	}

	return false, nil
}

func (p Playlist) AddToDatabase() error {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}

	var playlists []Playlist

	json.Unmarshal(byteValue, &playlists)

	log.Info("adding channel to DB")
	playlists = append(playlists, p)
	err = writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
	if err != nil {
		return fmt.Errorf("AddToDatabase: %s", err)
	}
	return nil
}

func (p Playlist) UpdateLatestDownloaded(videoID string) error {
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
		if playlist.PlaylistURL == p.PlaylistURL {
			playlists[i].LatestDownloaded = videoID
			log.Info("latest downloaded video id updated successfully")
			break
		}
	}

	return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
}

func (p Playlist) UpdateDownloadHistory(videoID string) error {
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
		if playlist.PlaylistURL == p.PlaylistURL {
			playlists[i].DownloadHistory = append(p.DownloadHistory, videoID)
			log.Info(playlists)
			break
		}
	}

	return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
}

func (p Playlist) GetFromDatabase() (Playlist, error) {
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
	if err != nil {
		return Playlist{}, fmt.Errorf("p.GetFromDatabase: %s", err)
	}
	var playlists []Playlist
	json.Unmarshal(byteValue, &playlists)

	for _, playlist := range playlists {
		if playlist.PlaylistURL == p.PlaylistURL {
			return playlist, nil
		}
	}

	return Playlist{}, fmt.Errorf("Couldn't find playlist in the database: %s", p.PlaylistURL)
}

func (p Playlist) UpdateLastChecked() error {
	log.Info("updating last checked date and time for: ", p.Name)

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
		if playlist.PlaylistURL == p.PlaylistURL {
			dt := time.Now()
			playlists[i].LastChecked = dt.Format("01-02-2006 15:04:05")
			log.Info("last checked date and time updated successfully")
			break
		}
	}

	return writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
}

func (p Playlist) Delete() error {
	log.Info("Removing channel from database")
	byteValue, err := openJSONDatabase(CONFIG_ROOT + "playlists.json")
	if err != nil {
		return fmt.Errorf("p.Delete: %s", err)
	}
	var playlists []Playlist

	json.Unmarshal(byteValue, &playlists)

	for i, playlist := range playlists {
		if playlist.PlaylistURL == p.PlaylistURL {
			playlists = RemoveAtIndexPlaylist(playlists, i)
			log.Info("successfully removed playlist from playlists.json")
		}
	}

	writeToPlaylistsDb(playlists, CONFIG_ROOT+"playlists.json")
	return nil
}

func writeToPlaylistsDb(db []Playlist, dbName string) error {
	result, err := json.Marshal(db)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("writeToPlaylistsDb: %s", err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	err = ioutil.WriteFile(dbName, file, 0644)
	if err != nil {
		log.Error("There was an error writing to database: ", err)
		return fmt.Errorf("writeToPlaylistsDb: %s", err)
	}

	return nil
}
