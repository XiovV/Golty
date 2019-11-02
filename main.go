package main

import (
	"fmt"
	"strings"
)

func main() {
	channelURL := "https://www.youtube.com/user/NewRetroWave"

	channelName := strings.Split(channelURL, "/")[4]
	channelType := strings.Split(channelURL, "/")[3]

	if channelType == "user" {
		uploadsId := GetUserUploadsID(channelName)
		videoId, videoTitle := GetUserVideoData(uploadsId)

		DownloadVideoAndAudio(videoId, videoTitle)
	} else if channelType == "channel" {
		videoId, videoTitle := GetChannelVideoData(channelName)
		fmt.Println(videoId, videoTitle)
	}

}
