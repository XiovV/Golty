package main

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Download downloads a the latest video based on downloadMode
func (c Channel) Download(downloadMode, fileExtension, downloadQuality string) error {
	video, err := c.GetLatestVideo()
	if err != nil {
		log.Error("Download: ", err)
		return fmt.Errorf("Download: %s", err)
	}
	err = video.Download(downloadMode, fileExtension, downloadQuality)
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

func (c Channel) DownloadEntire() error {
	fileExtension := strings.Replace(c.PreferredExtensionForAudio, ".", "", 1)
	var downloadSetting string
	if c.DownloadMode == "Audio Only" {
		downloadSetting = "bestaudio[ext=" + fileExtension + "]"
	} else if c.DownloadMode == "Video And Audio" {
		downloadSetting = "bestvideo[ext=" + fileExtension + "]"
	}
	cmd := exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", c.ChannelURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return fmt.Errorf("DownloadEntire: %s | %s", err, out)
	}

	return fmt.Errorf("DownloadEntire: download mode cannot be nil")
}
