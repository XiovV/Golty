package main

import (
	"fmt"
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

	t.Execute(w, GetChannels())
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
	channelURL := r.FormValue("channelURL")
	fmt.Println(channelURL)
	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	channel := []string{channelName}

	CheckNow(channel, channelType)
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")
	UpdateChannelsDatabase(channelURL)

	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	if channelType == "user" {
		fmt.Println("USER")
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle)
		}
	} else if channelType == "channel" {
		fmt.Println("CHANNEL")
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		fmt.Println(videoId, videoTitle)
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
