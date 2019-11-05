package main

import (
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

	channels := GetChannels()
	t.Execute(w, Response{Channels: channels})
}

func HandleCheckChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	channelName, err := GetChannelName(channelURL)
	if err != nil {
		fmt.Println(err)
		ReturnResponse(w, err.Error())
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		fmt.Println(err)
		ReturnResponse(w, err.Error())
	}

	CheckNow(channelName, channelType)
}

func HandleCheckAll(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	// ReturnResponse(w, "Checking")

	CheckNow("", "")
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	downloadMode := r.FormValue("mode")
	UpdateChannelsDatabase(channelURL)

	channelName, err := GetChannelName(channelURL)
	if err != nil {
		fmt.Println(err)
		ReturnResponse(w, err.Error())
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		fmt.Println(err)
		ReturnResponse(w, err.Error())
	}

	Download(channelName, channelType, downloadMode)

	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
}
