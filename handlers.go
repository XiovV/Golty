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
	log.Info("received a request to add a channel")
	var channelData AddChannelPayload
	err := json.NewDecoder(r.Body).Decode(&channelData)
	if err != nil {
		log.Error(err, r.Body)
	}
	channel := Channel{ChannelURL: channelData.ChannelURL}
	channelMetadata := channel.GetMetadata()

	if channelData.DownloadMode == "Audio Only" {
		channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForAudio: channelData.FileExtension}
	} else if channelData.DownloadMode == "Video And Audio" {
		channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForVideo: channelData.FileExtension}
	}

	doesChannelExist := channel.DoesExist()
	if doesChannelExist == true {
		log.Info("This channel already exists")
		res := Response{Type: "Success", Key: "CHANNEL_ALREADY_EXISTS", Message: "This channel already exists"}
		json.NewEncoder(w).Encode(res)
	} else {
		channel.AddToDatabase()

		err := channel.Download(channelData.DownloadMode, channelData.FileExtension, channelData.DownloadQuality)
		if err != nil {
			log.Error(err)
		}
		res := Response{Type: "Success", Key: "ADD_CHANNEL_SUCCESS", Message: "Channel successfully added"}
		json.NewEncoder(w).Encode(res)
	}
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check a channel for new uploads")
	w.Header().Set("Content-Type", "application/json")
	var data AddChannelPayload
	_ = json.NewDecoder(r.Body).Decode(&data)
	channel := Channel{ChannelURL: data.ChannelURL}

	log.Info(data)

	res := channel.CheckNow()
	json.NewEncoder(w).Encode(res)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check all channels for new uploads")
	w.Header().Set("Content-Type", "application/json")
	res := CheckAll()
	json.NewEncoder(w).Encode(res)
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to get all channels")
	w.Header().Set("Content-Type", "application/json")

	channels := GetChannels()

	json.NewEncoder(w).Encode(channels)
}

func HandleDeleteChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to delete a channel")

	w.Header().Set("Content-Type", "application/json")
	var data Payload
	_ = json.NewDecoder(r.Body).Decode(&data)
	channelURL := data.ChannelURL
	channelURL = strings.Replace(channelURL, "delChannel", "", -1)
	channel := Channel{ChannelURL: channelURL}

	channel.Delete()

	json.NewEncoder(w).Encode(Response{Type: "Success", Message: "Channel removed"})
}
