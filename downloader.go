package main

import (
	"fmt"
)

// Download downloads a the latest video based on downloadMode
func (c ChannelBasicInfo) Download(downloadMode, fileExtension, downloadQuality string) error {
	channelName := c.Name
	if downloadMode == "Video And Audio" {
		// Download .mp4 with audio and video in one file
		video := c.GetLatestVideo()
		return video.downloadVideoAndAudio(channelName, fileExtension, downloadQuality)
	} else if downloadMode == "Audio Only" {
		// Extract audio from the .mp4 file and remove the .mp4
		video := c.GetLatestVideo()
		return video.downloadAudioOnly(channelName, fileExtension, downloadQuality)
	}
	return fmt.Errorf("From Download: Something went seriously wrong")
}

func (v Video) downloadVideoAndAudio(channelName, fileExtension, downloadQuality string) error {
	err := v.DownloadYTDL(fileExtension, downloadQuality)
	if err != nil {
		return err
	}
	return UpdateLatestDownloaded(channelName, v.VideoID)
}

func (v Video) downloadAudioOnly(channelName, fileExtension, downloadQuality string) error {
	err := v.DownloadAudioYTDL(fileExtension, downloadQuality)
	if err != nil {
		return err
	}
	return UpdateLatestDownloaded(channelName, v.VideoID)
}
