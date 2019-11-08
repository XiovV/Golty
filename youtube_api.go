package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetLatestVideo(channelName, channelType string) (string, string) {
	if channelType == "user" {
		// Read the uploads id from a database before making an api call
		uploadsId, exists := GetUploadsIDFromDatabase(channelName)

		if exists == false {
			uploadsId = GetUserUploadsID(channelName)
			UpdateUploadsID(channelName, uploadsId)
		}

		requestURL := API_ENDPOINT_PLAYLIST + uploadsId + "&maxResults=2" + "&key=" + API_KEY

		resp, err := http.Get(requestURL)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		var video NameBody

		json.Unmarshal([]byte(string(body)), &video)

		return video.Items[0].Snippet.ResourceID.VideoID, video.Items[0].Snippet.Title
	}
	requestURL := API_ENDPOINT_ID + channelName + "&" + MAX_RESULTS + "&" + ORDER_BY + "&" + TYPE + "&key=" + API_KEY

	resp, err := http.Get(requestURL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var video ChannelBody

	json.Unmarshal([]byte(string(body)), &video)

	return video.Items[0].ID.VideoID, video.Items[0].Snippet.Title
}

func GetUserUploadsID(channelName string) string {
	fmt.Println("MAKING API CALL TO GET UPLOADS ID")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	requestURL := API_ENDPOINT_NAME + channelName + "&key=" + API_KEY

	resp, err := http.Get(requestURL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var user UserBody

	json.Unmarshal([]byte(string(body)), &user)

	return user.Items[0].ContentDetails.RelatedPlaylists.Uploads
}
