package main

import (
	"fmt"
)

// Download downloads a video based on downloadMode
func Download(channelName, channelType, downloadMode string) error {
	if downloadMode == "Video And Audio" {
		// Download .mp4 with audio and video in one file
		return downloadVideoAndAudio(channelName, channelType)
		// Extract audio from the .mp4 file and remove the .mp4
	} else if downloadMode == "Audio Only" {
		return downloadAudioOnly(channelName, channelType)
	}
	return fmt.Errorf("Something went seriously wrong")
}

func downloadVideoAndAudio(channelName, channelType string) error {
	videoId := GetLatestVideo(channelName, channelType)
	err := DownloadYTDL(videoId)
	if err != nil {
		return err
	}

	return UpdateLatestDownloaded(channelName, videoId)
}

func downloadAudioOnly(channelName, channelType string) error {
	videoId := GetLatestVideo(channelName, channelType)
	err := DownloadAudioYTDL(videoId)
	if err != nil {
		return err
	}

	return UpdateLatestDownloaded(channelName, videoId)
}
