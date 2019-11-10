package main

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func ServeJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/app.js")
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		log.Error(err)
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		log.Error(err)
	}

	doesChannelExist := DoesChannelExist(channelURL)
	if doesChannelExist == true {
		log.Info("This channel already exists")
	} else {
		// If the directory of the channel doesn't exist on the filesystem, create it
		CreateDirIfNotExist(channelName)
		log.Info("Adding channel to DB")
		AddChannelToDatabase(channelURL)
		if channelType == "user" {
			uploadsId := GetUserUploadsIDFromAPI(channelName)
			InitUploadsID(channelURL)
			UpdateUploadsID(channelURL, uploadsId)
			Download(channelName, channelType, downloadMode)
		} else if channelType == "channel" {
			Download(channelName, channelType, downloadMode)
		}
	}
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var channel Payload
	_ = json.NewDecoder(r.Body).Decode(&channel)
	channelURL := channel.ChannelURL

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		log.Error("There was an error getting the channel name: ", err)
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		log.Error("There was an error getting the channel type: ", err)
	}

	res := CheckNow(channelName, channelType)
	json.NewEncoder(w).Encode(res)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := CheckAll()
	json.NewEncoder(w).Encode(res)
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	channels := GetChannels()

	json.NewEncoder(w).Encode(channels)
}

func HandleDeleteChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var channel Payload
	_ = json.NewDecoder(r.Body).Decode(&channel)
	channelURL := channel.ChannelURL
	channelURL = strings.Replace(channelURL, "delChannel", "", -1)

	DeleteChannel(channelURL)

	json.NewEncoder(w).Encode(Response{Type: "Success", Message: "Channel removed"})
}
