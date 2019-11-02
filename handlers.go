package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")
	UpdateChannelsDatabase(channelURL)

	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	if channelType == "user" {
		fmt.Println("USER")
		uploadsId := GetUserUploadsID("NewRetroWave")
		videoId, videoTitle := GetUserVideoData(uploadsId)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle)
		}
	} else if channelType == "channel" {
		fmt.Println("CHANNEL")
		videoId, videoTitle := GetChannelVideoData(channelName)
		fmt.Println(videoId, videoTitle)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle)
		}
	}

	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}
