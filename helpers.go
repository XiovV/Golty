package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetChannels returns the contents of channels.json
func GetChannels() []Channel {
	log.Info("getting all channels from channels.json")
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
	log.Info("successfully read all channels")
	return db
}

// CheckAll goes through channels.json and checks for new videos
func CheckAll() Response {
	log.Info("checking for all channels")
	allChannelsInDb := GetChannels()
	var foundFor []string
	var preferredExtension string

	for _, item := range allChannelsInDb {
		channel := Channel{ChannelURL: item.ChannelURL}
		channel = channel.GetFromDatabase()
		// TODO: Handle errors

		if item.ChannelURL == channel.ChannelURL {
			videoId := channel.GetLatestVideo()

			if item.LatestDownloaded == videoId.VideoID {
				log.Info("no new videos found for: ", item.ChannelURL)
			} else {
				log.Info("new video detected for: ", item.ChannelURL)
				foundFor = append(foundFor, item.ChannelURL)
				if channel.DownloadMode == "Audio Only" {
					preferredExtension = channel.PreferredExtensionForAudio
				} else if channel.DownloadMode == "Video And Audio" {
					preferredExtension = channel.PreferredExtensionForVideo
				}
				go channel.Download(channel.DownloadMode, preferredExtension, "best")
				channel.UpdateLatestDownloaded(videoId.VideoID)
			}
		}
	}

	return Response{Type: "Success", Key: "NEW_VIDEOS_FOR_CHANNELS", Message: strings.Join(foundFor, ",")}
}

// CheckNow requires c.ChannelURL
func (c Channel) CheckNow() Response {
	log.Info("checking for new videos")
	allChannelsInDb := GetChannels()

	var preferredExtension string

	channel := c.GetFromDatabase()
	channelURL := c.ChannelURL

	channelMetadata := channel.GetMetadata()

	for _, item := range allChannelsInDb {
		if item.ChannelURL == channelURL {
			if item.LatestDownloaded == channelMetadata.ID {
				log.Info("no new videos found for: ", channelURL)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected"}
			} else {
				log.Info("new video detected for: ", channelURL)
				if channel.DownloadMode == "Audio Only" {
					preferredExtension = channel.PreferredExtensionForAudio
				} else if channel.DownloadMode == "Video And Audio" {
					preferredExtension = channel.PreferredExtensionForVideo
				}
				err := channel.Download(channel.DownloadMode, preferredExtension, "best")
				if err != nil {
					log.Error(err)
					return Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: err.Error()}
				}
				channel.UpdateLatestDownloaded(channelMetadata.ID)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected"}
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}
}

func CreateDirIfNotExist(dirName string) {
	log.Info("creating channel directory")
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
