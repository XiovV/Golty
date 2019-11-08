package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var channel Payload
	_ = json.NewDecoder(r.Body).Decode(&channel)
	channelURL := channel.ChannelURL

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		fmt.Println(err)
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CHECKING")

	CheckNow(channelName, channelType)

	json.NewEncoder(w).Encode(Response{Type: "Success", Message: "New videos detected"})

}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	go CheckAll()
	json.NewEncoder(w).Encode(Response{Type: "In Progress", Message: "Checking all channels"})
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	channelURL := r.FormValue("channelURL")
	// downloadMode := r.FormValue("mode")

	channelExists := UpdateChannelsDatabase(channelURL)
	UpdateUploadsIDDatabase(channelURL)

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		fmt.Println(err)
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		fmt.Println(err)
	}

	// Put uploads id of the user into a database
	if channelType == "user" {
		uploadsId := GetUserUploadsID(channelName)
		UpdateUploadsID(channelName, uploadsId)
	}

	// If the directory of the channel doesn't exist on the filesystem, create it
	CreateDirIfNotExist(channelName)

	if channelExists == false {
		DownloadVideoAndAudio(channelName, channelType)
	}
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	channels := GetChannels()

	json.NewEncoder(w).Encode(channels)
}
