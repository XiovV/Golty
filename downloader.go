package main

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Download downloads a the latest video based on downloadMode
func (c Channel) Download(downloadMode, fileExtension, downloadQuality string) error {
	channelURL := c.ChannelURL
	if downloadMode == "Video And Audio" {
		// Download .mp4 with audio and video in one file
		video := c.GetLatestVideo()
		return video.downloadVideoAndAudio(channelURL, fileExtension, downloadQuality)
	} else if downloadMode == "Audio Only" {
		// Extract audio from the .mp4 file and remove the .mp4
		video := c.GetLatestVideo()
		return video.downloadAudioOnly(channelURL, fileExtension, downloadQuality)
	}
	return fmt.Errorf("From Download: Something went seriously wrong")
}

func (c Channel) DownloadEntire() {
	if c.DownloadMode == "Audio Only" {
		fileExtension := strings.Replace(c.PreferredExtensionForAudio, ".", "", 1)
		cmd := exec.Command("youtube-dl", "-f", "bestaudio[ext="+fileExtension+"]", "-o", "downloads/ %(uploader)s/video/ %(title)s.%(ext)s", c.ChannelURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
		}
	} else if c.DownloadMode == "Video And Audio" {
		fileExtension := strings.Replace(c.PreferredExtensionForVideo, ".", "", 1)

		cmd := exec.Command("youtube-dl", "-f", "bestvideo[ext="+fileExtension+"]", "-o", "downloads/ %(uploader)s/video/ %(title)s.%(ext)s", c.ChannelURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
		}
	}

}

func (v Video) downloadVideoAndAudio(channelURL, fileExtension, downloadQuality string) error {
	err := v.DownloadYTDL(fileExtension, downloadQuality)
	if err != nil {
		return err
	}
	channel := Channel{ChannelURL: channelURL}
	return channel.UpdateLatestDownloaded(v.VideoID)
}

func (v Video) downloadAudioOnly(channelURL, fileExtension, downloadQuality string) error {
	err := v.DownloadAudioYTDL(fileExtension, downloadQuality)
	if err != nil {
		return err
	}
	channel := Channel{ChannelURL: channelURL}
	return channel.UpdateLatestDownloaded(v.VideoID)
}
