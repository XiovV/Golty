package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetChannels() []Channel {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	return db
}

func UploadChecker() {
	for {
		time.Sleep(10 * time.Second)
		// go CheckNow(nil, "")

		fmt.Println("Upload Checker running...")
	}
}

func CheckNow(channel string, channelType string) {
	allChannelsInDb := GetChannels()

	// if channel and channelType are both 0 then that means we want to check
	// for new uploads for all channels in the database
	if channel == "" && channelType == "" {
		for _, item := range allChannelsInDb {
			channelName, err := GetChannelName(item.ChannelURL)
			if err != nil {
				fmt.Println(err)
			}
			channelType, err := GetChannelType(item.ChannelURL)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(channelName, channelType)

			if strings.Contains(item.ChannelURL, channelName) {
				videoId, videoTitle := GetLatestVideo(channelName, channelType)

				if item.LatestDownloaded == videoId {
					fmt.Println("NOW NEW VIDEOS DETECTED FOR: ", item.ChannelURL)
				} else {
					fmt.Println("DOWNLOAD FOR: ", item.ChannelURL)
					DownloadAudio(videoId, videoTitle, channel)
					UpdateLatestDownloaded(item.ChannelURL, videoId)
				}
			}

		}
	} else {
		videoId, videoTitle := GetLatestVideo(channel, channelType)

		for _, item := range allChannelsInDb {
			if strings.Contains(item.ChannelURL, channel) {
				if item.LatestDownloaded == videoId {
					break
				} else {
					DownloadAudio(videoId, videoTitle, channel)
					UpdateLatestDownloaded(channel, videoId)
				}
			}
		}
	}
}

func ReturnResponse(w http.ResponseWriter, response string) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	channels := GetChannels()
	t.Execute(w, Response{Channels: channels, Status: response})
}

func GetChannelName(channelURL string) (string, error) {
	if channelURL != "" {
		return strings.Split(channelURL, "/")[4], nil
	}

	return "", fmt.Errorf("channelURL string is either empty or cant be parsed properly")
}

func GetChannelType(channelURL string) (string, error) {
	if channelURL != "" {
		return strings.Split(channelURL, "/")[3], nil
	}

	return "", fmt.Errorf("channelURL string is either empty or cant be parsed properly")
}

func CreateDirIfNotExist(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			panic(err)
		}
	}
}
