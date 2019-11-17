package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

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
