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
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	channel := Channel{ChannelURL: channelData.ChannelURL}

	doesChannelExist, err := channel.DoesExist()
	if err != nil {
		log.Info("error doesChannelExist: ", err)
		ReturnResponse(w, Response{Type: "Error", Key: "DOES_EXIST_ERROR", Message: "There was an error while trying to see if the channel already exists" + err.Error()})
	}
	if doesChannelExist == true {
		log.Info("this channel already exists")
		ReturnResponse(w, Response{Type: "Success", Key: "CHANNEL_ALREADY_EXISTS", Message: "This channel already exists"})
	} else {
		channelMetadata, err := channel.GetMetadata()
		if err != nil {
			ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting channel metadata: " + err.Error()})
		}

		if channelData.DownloadMode == "Audio Only" {
			channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForAudio: channelData.FileExtension}
		} else if channelData.DownloadMode == "Video And Audio" {
			channel = Channel{ChannelURL: channelData.ChannelURL, DownloadMode: channelData.DownloadMode, Name: channelMetadata.Uploader, PreferredExtensionForVideo: channelData.FileExtension}
		}
		err = channel.AddToDatabase()
		if err != nil {
			log.Error(err)
			ReturnResponse(w, Response{Type: "Error", Key: "ERROR_ADDING_CHANNEL", Message: "There was an error adding the channel to the database" + err.Error()})
		}
		if channelData.DownloadEntireChannel == true {
			err := channel.DownloadEntire()
			if err != nil {
				ReturnResponse(w, Response{Type: "Error", Key: "ERROR_DOWNLOADING_ENTIRE_CHANNEL", Message: "There was an error downloading the entire channel" + err.Error()})
			}
		} else {
			err = channel.AddToDatabase()
			if err != nil {
				log.Error(err)
				ReturnResponse(w, Response{Type: "Error", Key: "ERROR_ADDING_CHANNEL", Message: "There was an error adding the channel to the database" + err.Error()})
			}
			err = channel.Download(channelData.DownloadMode, channelData.FileExtension, channelData.DownloadQuality)
			if err != nil {
				log.Error(err)
				ReturnResponse(w, Response{Type: "Error", Key: "ERROR_DOWNLOADING", Message: "There was an error while downloading: " + err.Error()})
			}

			ReturnResponse(w, Response{Type: "Success", Key: "ADD_CHANNEL_SUCCESS", Message: "Channel successfully added and downloaded latest video"})
		}
	}
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check a channel for new uploads")
	w.Header().Set("Content-Type", "application/json")
	var data AddChannelPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	channel := Channel{ChannelURL: data.ChannelURL}

	res, err := channel.CheckNow()
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_CHECKING_CHANNEL", Message: "There was an error while checking the channel: " + err.Error()})
	}
	ReturnResponse(w, res)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check all channels for new uploads")
	w.Header().Set("Content-Type", "application/json")
	res, err := CheckAll()
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_CHECKING_CHANNELS", Message: "There was an error while checking channels: " + err.Error()})
	}
	ReturnResponse(w, res)
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to get all channels")
	w.Header().Set("Content-Type", "application/json")

	channels, err := GetChannels()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting channels: " + err.Error()}
		json.NewEncoder(w).Encode(res)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting channels: " + err.Error()})
	}
	json.NewEncoder(w).Encode(channels)
}

func HandleDeleteChannel(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to delete a channel")

	w.Header().Set("Content-Type", "application/json")
	var data Payload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	channelURL := data.ChannelURL
	channelURL = strings.Replace(channelURL, "delChannel", "", -1)
	channel := Channel{ChannelURL: channelURL}

	channel.Delete()

	ReturnResponse(w, Response{Type: "Success", Key: "DELETE_CHANNEL_SUCCESS", Message: "Channel removed"})
}
