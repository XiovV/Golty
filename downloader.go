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
		video, err := c.GetLatestVideo()
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		// video.downloadVideoAndAudio(channelURL, fileExtension, downloadQuality)
		err = video.DownloadVideoYTDL(fileExtension, downloadQuality)
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		err = c.UpdateLatestDownloaded(video.VideoID)
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		return c.UpdateDownloadHistory(video.VideoID)
	} else if downloadMode == "Audio Only" {
		// Extract audio from the .mp4 file and remove the .mp4
		video, err := c.GetLatestVideo()
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		// video.downloadAudioOnly(channelURL, fileExtension, downloadQuality)
		err = video.DownloadAudioYTDL(fileExtension, downloadQuality)
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		err = c.UpdateLatestDownloaded(video.VideoID)
		if err != nil {
			log.Error("Download: ", err)
			return fmt.Errorf("Download: %s", err)
		}
		return c.UpdateDownloadHistory(video.VideoID)
	}
	return fmt.Errorf("From Download: Something went seriously wrong")
}

func (c Channel) DownloadEntire() error {
	if c.DownloadMode == "Audio Only" {
		fileExtension := strings.Replace(c.PreferredExtensionForAudio, ".", "", 1)
		cmd := exec.Command("youtube-dl", "-f", "bestaudio[ext="+fileExtension+"]", "--ignore-errors", "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", c.ChannelURL)
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
