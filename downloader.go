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
		log.Error("c.Download: ", err)
		return fmt.Errorf("c.Download: %s", err)
	}
	err = video.Download(downloadMode, fileExtension, downloadQuality)
	if err != nil {
		log.Error("c.Download: ", err)
		return fmt.Errorf("c.Download: %s", err)
	}
	err = c.UpdateLatestDownloaded(video.VideoID)
	if err != nil {
		log.Error("c.Download: ", err)
		return fmt.Errorf("c.Download: %s", err)
	}
	return c.UpdateDownloadHistory(video.VideoID)
}

func (p Playlist) Download(downloadMode, fileExtension, downloadQuality string) error {
	video, err := p.GetLatestVideo()
	if err != nil {
		log.Error("p.Download: ", err)
		return fmt.Errorf("p.Download: %s", err)
	}
	err = video.Download(downloadMode, fileExtension, downloadQuality)
	if err != nil {
		log.Error("p.Download: ", err)
		return fmt.Errorf("p.Download: %s", err)
	}
	err = p.UpdateLatestDownloaded(video.VideoID)
	if err != nil {
		log.Error("p.Download: ", err)
		return fmt.Errorf("p.Download: %s", err)
	}
	return p.UpdateDownloadHistory(video.VideoID)
}

func (p Playlist) DownloadEntire() error {
	fileExtension := strings.Replace(p.PreferredExtensionForAudio, ".", "", 1)
	var cmd *exec.Cmd
	var downloadSetting string
	if p.DownloadMode == "Audio Only" {
		downloadSetting = "bestaudio[ext=" + fileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/playlists/%(uploader)s/audio/%(title)s.%(ext)s", p.PlaylistURL)
	} else if p.DownloadMode == "Video And Audio" {
		downloadSetting = "bestvideo[ext=" + fileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/playlists/%(uploader)s/video/%(title)s.%(ext)s", p.PlaylistURL)
	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("p.DownloadEntire: %s", err)
	}

	return fmt.Errorf("DownloadEntire: download mode cannot be nil")
}

func (c Channel) DownloadEntire() error {
	var cmd *exec.Cmd
	fileExtension := strings.Replace(c.PreferredExtensionForAudio, ".", "", 1)
	var downloadSetting string
	if c.DownloadMode == "Audio Only" {
		downloadSetting = "bestaudio[ext=" + fileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", c.ChannelURL)
	} else if c.DownloadMode == "Video And Audio" {
		downloadSetting = "bestvideo[ext=" + fileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/%(uploader)s/video/%(title)s.%(ext)s", c.ChannelURL)

	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("c.DownloadEntire: %s", err)
	}

	return fmt.Errorf("DownloadEntire: download mode cannot be nil")
}
