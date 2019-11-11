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

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "log.log")
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	var channelData AddChannelPayload
	err := json.NewDecoder(r.Body).Decode(&channelData)
	if err != nil {
		log.Error(err, r.Body)
	}
	channelURL := channelData.ChannelURL

	downloadMode := channelData.DownloadMode

	log.Info(channelData)

	channel, _ := GetChannelInfo(channelURL)

	channelType := channel.Type

	doesChannelExist := DoesChannelExist(channelURL)
	if doesChannelExist == true {
		log.Info("This channel already exists")
		res := Response{Type: "Success", Key: "CHANNEL_ALREADY_EXISTS", Message: "This channel already exists"}
		json.NewEncoder(w).Encode(res)
	} else {
		log.Info("Adding channel to DB")
		AddChannelToDatabase(channelURL)
		if channelType == "user" {
			err := channel.Download(downloadMode)
			if err != nil {
				log.Error(err)
			}
			res := Response{Type: "Success", Key: "ADD_CHANNEL_SUCCESS", Message: "Channel successfully added"}
			json.NewEncoder(w).Encode(res)
		} else if channelType == "channel" {
			err := channel.Download(downloadMode)
			if err != nil {
				log.Error(err)
			}
			res := Response{Type: "Success", Key: "ADD_CHANNEL_SUCCESS", Message: "Channel successfully added"}
			json.NewEncoder(w).Encode(res)
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
