package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetChannels() []Channel {
	log.Info("Getting all channels from channels.json")
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Error("There was an error reading channels.json: ", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		log.Error("There was an error unmarshalling json: ", err)
	}
	log.Info("Successfully read all channels")
	return db
}

func GetChannelInfo(channelURL string) (ChannelBasicInfo, error) {
	channelName, err := GetChannelName(channelURL)
	if err != nil {
		return ChannelBasicInfo{}, fmt.Errorf("There was an error getting channel name: %s", err)
	}
	channelType, err := GetChannelType(channelURL)
	if err != nil {
		return ChannelBasicInfo{}, fmt.Errorf("There was an error getting channel type: %s", err)
	}

	return ChannelBasicInfo{Name: channelName, Type: channelType}, nil
}

func CheckAll() Response {
	log.Info("Checking for all channels")
	allChannelsInDb := GetChannels()
	var foundFor []string

	for _, item := range allChannelsInDb {
		channel, err := GetChannelInfo(item.ChannelURL)
		if err != nil {
			log.Error(err)
		}
		channelName := channel.Name

		if strings.Contains(item.ChannelURL, channelName) {
			videoId := channel.GetLatestVideo()
			// videoId := GetLatestVideo(channelName, channelType)

			if item.LatestDownloaded == videoId.VideoID {
				log.Info("No new videos found for: ", item.ChannelURL)
			} else {
				log.Info("New video detected for: ", item.ChannelURL)
				foundFor = append(foundFor, item.ChannelURL)
				go channel.Download("Audio Only", ".mp3", "best")
				UpdateLatestDownloaded(item.ChannelURL, videoId.VideoID)
			}
		}
	}

	return Response{Type: "Success", Key: "NEW_VIDEOS_FOR_CHANNELS", Message: strings.Join(foundFor, ",")}
}

func CheckNow(channelName string, channelType string) Response {
	log.Info("Checking for new videos")
	allChannelsInDb := GetChannels()

	channel := ChannelBasicInfo{Name: channelName, Type: channelType}

	videoId := channel.GetLatestVideo()

	for _, item := range allChannelsInDb {
		if strings.Contains(item.ChannelURL, channelName) {
			if item.LatestDownloaded == videoId.VideoID {
				log.Info("No new videos found for: ", channelName)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected"}
			} else {
				log.Info("New video detected for: ", channelName)
				err := channel.Download("Audio Only", ".mp3", "best")
				if err != nil {
					log.Error(err)
					return Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: err.Error()}
				}
				UpdateLatestDownloaded(channelName, videoId.VideoID)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected"}
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}
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
	log.Info("Creating channel directory")
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			log.Error("Couldn't create channel directory: ", err)
		} else {
			log.Info("Channel directory created successfully")
		}
	}
}

func RemoveAtIndex(s []Channel, index int) []Channel {
	return append(s[:index], s[index+1:]...)
}

func GetFailedDownloads() []FailedVideo {
	log.Info("Getting failed downloads")
	jsonFile, err := os.Open("failed.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []FailedVideo

	json.Unmarshal(byteValue, &db)

	return db
}
