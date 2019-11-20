package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//type AddChannelPayload struct {
//	URL            string
//	DownloadMode          string
//	FileExtension         string
//	DownloadQuality       string
//	DownloadEntire bool
//	CheckingInterval      string
//}

func readDownloadsChannelsDir() []string {
	var channels []string
	file, err := os.Open("./downloads/channels")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()



	list,_ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		channels = append(channels, name)
	}

	return channels
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/add-channel", HandleAddChannel).Methods("POST")
	return router
}

var testChannels = []AddChannelPayload{
	{
		URL:              "https://www.youtube.com/user/HungOverGargoyle",
		DownloadMode:     "Audio Only",
		FileExtension:    "m4a",
		DownloadQuality:  "best",
		DownloadEntire:   false,
		CheckingInterval: "",
	},
	{
		URL:              "https://www.youtube.com/user/NewRetroWave",
		DownloadMode:     "Video And Audio",
		FileExtension:    "mp4",
		DownloadQuality:  "best",
		DownloadEntire:   false,
		CheckingInterval: "",
	},
}

func TestHandleAddChannelWithChannelsInDatabase(t *testing.T) {
	// 1. Add channel to db
	// 2. Make request to /api/add-channel/ to see if it's going to return a response saying that the channel already exists
	assert := assert.New(t)
	var res Response
	for _, testChannel := range testChannels {
		// Adding a channel into db
		channel := DownloadTarget{URL: testChannel.URL, Type:"Channel"}
		err := channel.AddToDatabase()
		assert.Nil(err)
		// Marshalling the json object and making POST request
		channelJson, _ := json.Marshal(testChannel)
		request, _ := http.NewRequest("POST", "/api/add-channel", bytes.NewBuffer(channelJson))
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
		err = json.NewDecoder(response.Body).Decode(&res)
		assert.Nil(err)
		// Checking the response
		assert.Equal("Success", res.Type, "res.Type should be success")
		assert.Equal("CHANNEL_ALREADY_EXISTS", res.Key)

		err = channel.Delete()
		assert.Nil(err)
	}
}

func TestHandleAddChannelWithoutAnyChannelsInDatabase(t *testing.T) {
	assert := assert.New(t)
	var res Response
	for _, testChannel := range testChannels {
		channelJson, _ := json.Marshal(testChannel)
		request, _ := http.NewRequest("POST", "/api/add-channel", bytes.NewBuffer(channelJson))
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
		err := json.NewDecoder(response.Body).Decode(&res)
		channel := DownloadTarget{URL: testChannel.URL, Type:"Channel"}
		assert.Nil(err)
		assert.Equal("Success", res.Type, "res.Type should be success")
		assert.Equal("ADD_CHANNEL_SUCCESS", res.Key)
		addedChannel, err := channel.GetFromDatabase()
		assert.Nil(err)
		if testChannel.URL == testChannels[0].URL {
			assert.Equal("https://www.youtube.com/user/HungOverGargoyle", addedChannel.URL)
			assert.Equal("Audio Only", addedChannel.DownloadMode)
			assert.Equal("m4a", addedChannel.PreferredExtensionForAudio)
			assert.Equal(1, len(addedChannel.DownloadHistory))
		} else if testChannel.URL == testChannels[1].URL {
			assert.Equal("https://www.youtube.com/user/NewRetroWave", addedChannel.URL)
			assert.Equal("Video And Audio", addedChannel.DownloadMode)
			assert.Equal("mp4", addedChannel.PreferredExtensionForVideo)
			assert.Equal(1, len(addedChannel.DownloadHistory))
		}
		err = channel.Delete()
		assert.Nil(err)
	}
	downloadedChannels := readDownloadsChannelsDir()
	fmt.Println("downloadedChannels: ", downloadedChannels)

	assert.Equal(2, len(downloadedChannels))

	err := os.RemoveAll("./downloads")
	assert.Nil(err)
}
