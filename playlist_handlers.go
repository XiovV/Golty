package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func HandleAddPlaylist(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to add a playlist")
	w.Header().Set("Content-Type", "application/json")
	var data AddPlaylistPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
}
