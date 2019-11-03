package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	channels := GetChannels()
	t.Execute(w, Response{Channels: channels})
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	channel := []string{channelName}

	CheckNow(channel, channelType)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	ReturnResponse(w, "Checking")

	CheckNow(nil, "")
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")
	UpdateChannelsDatabase(channelURL)

	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	if channelType == "user" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle)
		}
	} else if channelType == "channel" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle)
		}
	}

	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
}
