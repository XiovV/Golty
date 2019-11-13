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
		res := Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}
	channel := Channel{ChannelURL: channelData.ChannelURL}
	channelMetadata, err := channel.GetMetadata()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting channel metadata: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}

	if channelData.DownloadMode == "Audio Only" {
		channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForAudio: channelData.FileExtension}
	} else if channelData.DownloadMode == "Video And Audio" {
		channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForVideo: channelData.FileExtension}
	}

	doesChannelExist := channel.DoesExist()
	if doesChannelExist == true {
		log.Info("this channel already exists")
		res := Response{Type: "Success", Key: "CHANNEL_ALREADY_EXISTS", Message: "This channel already exists"}
		json.NewEncoder(w).Encode(res)
	} else {
		channel.AddToDatabase()

		err := channel.Download(channelData.DownloadMode, channelData.FileExtension, channelData.DownloadQuality)
		if err != nil {
			log.Error(err)
			res := Response{Type: "Error", Key: "ERROR_DOWNLOADING", Message: "There was an error while downloading: " + err.Error()}
			json.NewEncoder(w).Encode(res)
		}
		res := Response{Type: "Success", Key: "ADD_CHANNEL_SUCCESS", Message: "Channel successfully added"}
		json.NewEncoder(w).Encode(res)
	}
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check a channel for new uploads")
	w.Header().Set("Content-Type", "application/json")
	var data AddChannelPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}
	channel := Channel{ChannelURL: data.ChannelURL}

	res, err := channel.CheckNow()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_CHECKING_CHANNEL", Message: "There was an error while checking the channel: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}
	json.NewEncoder(w).Encode(res)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check all channels for new uploads")
	w.Header().Set("Content-Type", "application/json")
	res, err := CheckAll()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_CHECKING_CHANNELS", Message: "There was an error while checking channels: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}
	json.NewEncoder(w).Encode(res)
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to get all channels")
	w.Header().Set("Content-Type", "application/json")

	channels, err := GetChannels()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting channels: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}

	json.NewEncoder(w).Encode(channels)
}

func HandleDeleteChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to delete a channel")

	w.Header().Set("Content-Type", "application/json")
	var data Payload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
		json.NewEncoder(w).Encode(res)
	}
	channelURL := data.ChannelURL
	channelURL = strings.Replace(channelURL, "delChannel", "", -1)
	channel := Channel{ChannelURL: channelURL}

	channel.Delete()

	json.NewEncoder(w).Encode(Response{Type: "Success", Key: "DELETE_CHANNEL_SUCCESS", Message: "Channel removed"})
}
