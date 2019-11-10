package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetLatestVideo(channelName, channelType string) (string, string) {
	if channelType == "user" {
		// Read the uploads id from a database before making an api call
		uploadsId, doesUserExist := GetUploadsIDFromDatabase(channelName)

		if doesUserExist == false { // If the user doesn't exist then get uploadsId from the api and update it in the database
			log.Warn("User doesn't exist in uploadid.json, making api request to get uploads id")
			uploadsId = GetUserUploadsIDFromAPI(channelName)
			UpdateUploadsID(channelName, uploadsId)
		} else {
			log.Info("Found uploads id for user, not making an api request")
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

func GetUserUploadsIDFromAPI(channelName string) string {
	log.Println("Making api call to get uploads id for a user")
	requestURL := API_ENDPOINT_NAME + channelName + "&key=" + API_KEY

	resp, err := http.Get(requestURL)

	if err != nil {
		log.Error("Couldn't make api request", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("Couldn't read response body")
	}

	var user UserBody

	json.Unmarshal([]byte(string(body)), &user)

	return user.Items[0].ContentDetails.RelatedPlaylists.Uploads
}
