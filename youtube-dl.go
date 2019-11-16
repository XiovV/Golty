package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// GetMetadata only requires c.ChannelURL, it returns ChannelMetadata{} containing all metadata for the channel.
func (c Channel) GetMetadata() (ChannelMetadata, error) {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", c.ChannelURL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("From GetMetadata(): ", err)
		return ChannelMetadata{}, fmt.Errorf("From GetMetadata(): %v", err)
	}
	metaData := &ChannelMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("From GetMetadata(): ", err)
		return ChannelMetadata{}, fmt.Errorf("From GetMetadata(): %v", err)
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
		log.Errorf("GetLatestVideo: %s | %s", err, string(out))
	}
	metaData := &ChannelMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("GetLatestVideo: ", err)
		return Video{}, fmt.Errorf("GetLatestVideo: %s", err)
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

// // DownloadYTDL downloads a video file with specified paramaters
// func (v Video) DownloadVideoYTDL(fileExtension, downloadQuality string) error {
// 	log.Info("downloading video file")
// 	cmd := exec.Command("youtube-dl", "-f", downloadQuality, "-o", "downloads/%(uploader)s/video/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
// 	log.Info("executing youtube-dl command: ", cmd.String())

// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Fatal(string(out))
// 		return fmt.Errorf("DownloadVideoYTDL: %s", err)
// 	}

// 	log.Info("successfully downloaded video file")
// 	return nil
// }

// // DownloadAudioYTDL downloads an audio file with specified paramaters
// func (v Video) DownloadAudioYTDL(fileExtension, downloadQuality string) error {
// 	log.Info("downloading audio only")
// 	if downloadQuality == "best" {
// 		downloadQuality = "0"
// 	} else if downloadQuality == "medium" {
// 		downloadQuality = "5"
// 	} else if downloadQuality == "worst" {
// 		downloadQuality = "9"
// 	}
// 	log.Info("download quality set to: ", downloadQuality)
// 	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", fileExtension, "--audio-quality", downloadQuality, "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
// 	log.Info("executing youtube-dl command: ", cmd.String())
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Error(string(out))
// 		return fmt.Errorf("DownloadAudioYTDL: %s", err)
// 	}

// 	log.Info("successfully downloaded audio")
// 	return nil
// }
