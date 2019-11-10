package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func CheckAll() {
	allChannelsInDb := GetChannels()

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
			videoId, _ := GetLatestVideo(channelName, channelType)

			if item.LatestDownloaded == videoId {
				fmt.Println("NOW NEW VIDEOS DETECTED FOR: ", item.ChannelURL)
			} else {
				fmt.Println("DOWNLOAD FOR: ", item.ChannelURL)
				go DownloadAudio(channelName, channelType)
				UpdateLatestDownloaded(item.ChannelURL, videoId)
			}
		}
	}
}

func CheckNow(channel string, channelType string) Response {
	allChannelsInDb := GetChannels()

	videoId, _ := GetLatestVideo(channel, channelType)

	for _, item := range allChannelsInDb {
		if strings.Contains(item.ChannelURL, channel) {
			if item.LatestDownloaded == videoId {
				fmt.Println("No new videos")
				return Response{Type: "False", Message: "No new videos detected"}
			} else {
				DownloadAudio(channel, channelType)
				UpdateLatestDownloaded(channel, videoId)
				return Response{Type: "True", Message: "New video detected"}
			}
		}
	}
	return Response{Type: "Error", Message: "Something went wrong"}
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

func RemoveAtIndex(s []Channel, index int) []Channel {
	return append(s[:index], s[index+1:]...)
}
