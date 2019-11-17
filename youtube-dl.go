package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetMetadata only requires c.ChannelURL, it returns ChannelMetadata{} containing all metadata for the channel.
func (c Channel) GetMetadata() (ChannelMetadata, error) {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", c.ChannelURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("From GetMetadata(): ", err)
		return ChannelMetadata{}, fmt.Errorf("From c.GetMetadata(): %v", err)
	}
	metaData := &ChannelMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("From GetMetadata(): ", err)
		return ChannelMetadata{}, fmt.Errorf("From c.GetMetadata(): %v", err)
	}

	return *metaData, nil
}

func (p Playlist) GetMetadata() (PlaylistMetadata, error) {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", p.PlaylistURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("From GetMetadata(): ", err)
		return PlaylistMetadata{}, fmt.Errorf("From p.GetMetadata(): %v", err)
	}
	metaData := &PlaylistMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("From GetMetadata(): ", err)
		return PlaylistMetadata{}, fmt.Errorf("From p.GetMetadata(): %v", err)
	}

	return *metaData, nil
}

// GetLatestVideo only requires c.ChannelURL, it returns Video{} with VideoID
func (c Channel) GetLatestVideo() (Video, error) {
	log.Info("fetching latest upload")
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", c.ChannelURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		log.Errorf("c.GetLatestVideo: %s | %s", err, string(out))
	}
	metaData := &ChannelMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("c.GetLatestVideo: ", err)
		return Video{}, fmt.Errorf("c.GetLatestVideo: %s", err)
	}
	log.Info("successfully fetched latest video ")
	return Video{VideoID: metaData.ID}, nil
}

func (p Playlist) GetLatestVideo() (Video, error) {
	log.Info("fetching latest upload")
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", p.PlaylistURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		log.Errorf("p.GetLatestVideo: %s | %s", err, string(out))
	}
	metaData := &PlaylistMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("p.GetLatestVideo: ", err)
		return Video{}, fmt.Errorf("p.GetLatestVideo: %s", err)
	}
	log.Info("successfully fetched latest video ")
	return Video{VideoID: metaData.ID}, nil
}

func (v Video) Download(downloadMode, fileExtension, downloadQuality string) error {
	log.Info("executing download")
	var cmd *exec.Cmd
	if downloadMode == "Audio Only" {
		log.Info("downloading audio only")
		if downloadQuality == "best" {
			downloadQuality = "0"
		} else if downloadQuality == "medium" {
			downloadQuality = "5"
		} else if downloadQuality == "worst" {
			downloadQuality = "9"
		}
		log.Info("download quality set to: ", downloadQuality)
		cmd = exec.Command("youtube-dl", "--extract-audio", "--audio-format", fileExtension, "--audio-quality", downloadQuality, "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
	} else if downloadMode == "Video And Audio" {
		cmd = exec.Command("youtube-dl", "-f", downloadQuality, "-o", "downloads/%(uploader)s/video/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("v.Download: %s", err)
	}
	return nil
}

func (p Playlist) DownloadPlaylist(downloadMode, fileExtension, downloadQuality string) error {
	log.Info("executing download")
	var cmd *exec.Cmd
	if downloadMode == "Audio Only" {
		log.Info("downloading audio only")
		if downloadQuality == "best" {
			downloadQuality = "0"
		} else if downloadQuality == "medium" {
			downloadQuality = "5"
		} else if downloadQuality == "worst" {
			downloadQuality = "9"
		}
		log.Info("download quality set to: ", downloadQuality)
		cmd = exec.Command("youtube-dl", "--extract-audio", "--audio-format", fileExtension, "--playlist-end", "1", "--audio-quality", downloadQuality, "-o", "downloads/playlists/%(uploader)s/%(playlist)s/audio/%(title)s.%(ext)s", p.PlaylistURL)
	} else if downloadMode == "Video And Audio" {
		cmd = exec.Command("youtube-dl", "-f", downloadQuality, "--playlist-end", "1", "-o", "downloads/playlists/%(uploader)s/%(playlist)s/video/%(title)s.%(ext)s", p.PlaylistURL)
	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("v.DownloadPlaylist: %s", err)
	}
	return nil
}

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
	err = p.DownloadPlaylist(downloadMode, fileExtension, downloadQuality)
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

func (v DownloadVideoPayload) Download() error {
	var cmd *exec.Cmd
	var downloadSetting string
	log.Info(v.DownloadMode)
	if v.DownloadMode == "Audio Only" {
		downloadSetting = "bestaudio[ext=" + v.FileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/videos/%(uploader)s/audio/%(title)s.%(ext)s", v.VideoURL)
	} else if v.DownloadMode == "Video And Audio" {
		downloadSetting = "bestvideo[ext=" + v.FileExtension + "]"
		cmd = exec.Command("youtube-dl", "-f", downloadSetting, "--ignore-errors", "-o", "downloads/videos/%(uploader)s/video/%(title)s.%(ext)s", v.VideoURL)
	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("v.Download: %s", err)
	}

	return v.AddToDatabase()
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