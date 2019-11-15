package main

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Download downloads a the latest video based on downloadMode
func (c Channel) Download(downloadMode, fileExtension, downloadQuality string) error {
	if downloadMode == "Video And Audio" {
		// Download .mp4 with audio and video in one file
		video := c.GetLatestVideo()
		// video.downloadVideoAndAudio(channelURL, fileExtension, downloadQuality)
		video.DownloadVideoYTDL(fileExtension, downloadQuality)
		return c.UpdateLatestDownloaded(video.VideoID)
	} else if downloadMode == "Audio Only" {
		// Extract audio from the .mp4 file and remove the .mp4
		video := c.GetLatestVideo()
		// video.downloadAudioOnly(channelURL, fileExtension, downloadQuality)
		video.DownloadAudioYTDL(fileExtension, downloadQuality)
		return c.UpdateLatestDownloaded(video.VideoID)
	}
	return fmt.Errorf("From Download: Something went seriously wrong")
}

func (c Channel) DownloadEntire() error {
	if c.DownloadMode == "Audio Only" {
		fileExtension := strings.Replace(c.PreferredExtensionForAudio, ".", "", 1)
		cmd := exec.Command("youtube-dl", "-f", "bestaudio[ext="+fileExtension+"]", "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", c.ChannelURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return fmt.Errorf("DownloadEntire: %s | %s", err, out)
		}
	} else if c.DownloadMode == "Video And Audio" {
		fileExtension := strings.Replace(c.PreferredExtensionForVideo, ".", "", 1)

		cmd := exec.Command("youtube-dl", "-f", "bestvideo[ext="+fileExtension+"]", "-o", "downloads/%(uploader)s/video/%(title)s.%(ext)s", c.ChannelURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return fmt.Errorf("DownloadEntire: %s | %s", err, out)
		}
	}
	return fmt.Errorf("DownloadEntire: download mode cannot be nil")
}
