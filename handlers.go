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
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	// ReturnResponse(w, "Checking")

	CheckNow("", "")
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")

	exists := UpdateChannelsDatabase(channelURL)

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		fmt.Println(err)
		// ReturnResponse(w, err.Error())
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		fmt.Println(err)
		// ReturnResponse(w, err.Error())
	}

	// If the directory of the channel doesn't exist on the filesystem, create it
	CreateDirIfNotExist(channelName)

	if exists == false {
		Download(channelName, channelType, downloadMode)
		// ReturnResponse(w, "Channel added successfully")
	} else {
		// ReturnResponse(w, "This channel already exists")
	}
}

func HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	channels := GetChannels()

	json.NewEncoder(w).Encode(channels)
}
