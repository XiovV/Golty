package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func getFileNames(path string) (string, error){
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return files[0].Name(), nil
}

func checkVideoExtension(expectedExtension, pathToDirectory string) (bool, error) {
	videoName, err := getFileNames(pathToDirectory)
	if err != nil {
		return false, fmt.Errorf(err.Error())
	}
	videoExtension := strings.Split(videoName, ".")[1]
	if videoExtension != expectedExtension {
		return false, nil
	} else {
		return true, nil
	}
}

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

var testChannels = []AddTargetPayload{
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
			isExtensionCorrect, err := checkVideoExtension("m4a", "./downloads/channels/Electronic Gems/audio/")
			assert.Nil(err)
			assert.Equal(true, isExtensionCorrect)
		} else if testChannel.URL == testChannels[1].URL {
			assert.Equal("https://www.youtube.com/user/NewRetroWave", addedChannel.URL)
			assert.Equal("Video And Audio", addedChannel.DownloadMode)
			assert.Equal("mp4", addedChannel.PreferredExtensionForVideo)
			assert.Equal(1, len(addedChannel.DownloadHistory))
			isExtensionCorrect, err := checkVideoExtension("mkv", "./downloads/channels/NewRetroWave/video/")
			assert.Nil(err)
			assert.Equal(true, isExtensionCorrect)
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

//func TestHandleCheckChannelWithOutdatedLatestDownloaded(t *testing.T) {
//	assert := assert.New(t)
//	var res Response
//	for _, testChannel := range testChannels {
//		channelJson, _ := json.Marshal(testChannel)
//		channel := DownloadTarget{URL: testChannel.URL, Type: "Channel"}
//		channel.AddToDatabase()
//		channelJson, _ = json.Marshal(channel)
//		request, err := http.NewRequest("POST", "/api/check-channel", bytes.NewBuffer(channelJson))
//		assert.Nil(err)
//		response := httptest.NewRecorder()
//		Router().ServeHTTP(response, request)
//		fmt.Println(response.Code)
//		err = json.NewDecoder(response.Body).Decode(&res)
//		assert.Nil(err)
//		assert.Equal("Success", res.Type, "res.Type should be success")
//		assert.Equal("NEW_VIDEO_DETECTED", res.Key)
//		addedChannel, err := channel.GetFromDatabase()
//		assert.Nil(err)
//		if testChannel.URL == testChannels[0].URL {
//			fmt.Println("MADE API CALL 0")
//
//			assert.Equal("https://www.youtube.com/user/HungOverGargoyle", addedChannel.URL)
//			assert.Equal("Audio Only", addedChannel.DownloadMode)
//			assert.Equal("m4a", addedChannel.PreferredExtensionForAudio)
//			assert.Equal(1, len(addedChannel.DownloadHistory))
//			isExtensionCorrect, err := checkVideoExtension("m4a", "./downloads/channels/Electronic Gems/audio/")
//			assert.Nil(err)
//			assert.Equal(true, isExtensionCorrect)
//		} else if testChannel.URL == testChannels[1].URL {
//			fmt.Println("MADE API CALL 1")
//
//			assert.Equal("https://www.youtube.com/user/NewRetroWave", addedChannel.URL)
//			assert.Equal("Video And Audio", addedChannel.DownloadMode)
//			assert.Equal("mp4", addedChannel.PreferredExtensionForVideo)
//			assert.Equal(1, len(addedChannel.DownloadHistory))
//			isExtensionCorrect, err := checkVideoExtension("mkv", "./downloads/channels/NewRetroWave/video/")
//			assert.Nil(err)
//			assert.Equal(true, isExtensionCorrect)
//		}
//		err = channel.Delete()
//		assert.Nil(err)
//	}
//	downloadedChannels := readDownloadsChannelsDir()
//	fmt.Println("downloadedChannels: ", downloadedChannels)
//
//	assert.Equal(2, len(downloadedChannels))
//
//	err := os.RemoveAll("./downloads")
//	assert.Nil(err)
//}
