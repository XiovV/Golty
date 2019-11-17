package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func GetMetadata(target interface{}) (interface{}, error) {
	switch target := target.(type) {
		case Channel:
			cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.ChannelURL)
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

		case Playlist:
			cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.PlaylistURL)
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

	return nil, fmt.Errorf("target must be either Channel or Playlist")
}

func GetLatestVideo(target interface{}) (string, error){
	switch target := target.(type) {
	case Channel:
		log.Info("fetching latest upload")
		cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.ChannelURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			log.Errorf("c.GetLatestVideo: %s | %s", err, string(out))
		}
		metaData := &ChannelMetadata{}
		if err = json.Unmarshal(out, metaData); err != nil {
			log.Error("c.GetLatestVideo: ", err)
			return "", fmt.Errorf("c.GetLatestVideo: %s", err)
		}
		log.Info("successfully fetched latest video ")
		return metaData.ID, nil
	case Playlist:
		log.Info("fetching latest upload")
		cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.PlaylistURL)
		log.Info("executing youtube-dl command: ", cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			log.Errorf("p.GetLatestVideo: %s | %s", err, string(out))
		}
		metaData := &PlaylistMetadata{}
		if err = json.Unmarshal(out, metaData); err != nil {
			log.Error("p.GetLatestVideo: ", err)
			return "", fmt.Errorf("p.GetLatestVideo: %s", err)
		}
		log.Info("successfully fetched latest video ")
		fmt.Println(metaData.ID)
		return metaData.ID, nil
	}

	return "", fmt.Errorf("targetType has to be either Channel or Playlist")
}



func Download(target interface{}, downloadQuality, fileExtension string) error {
	switch target := target.(type) {
	case Channel:
		videoId, err := GetLatestVideo(target)
			if err != nil {
			log.Error("c.Download: ", err)
			return fmt.Errorf("c.Download: %s", err)
		}
		video := Video{VideoID: videoId, DownloadMode: target.DownloadMode}
		Download(video, downloadQuality, fileExtension)
		err = UpdateLatestDownloaded(target, video.VideoID)
		if err != nil {
			log.Error("p.Download: ", err)
			return fmt.Errorf("p.Download: %s", err)
		}
		return UpdateDownloadHistory(target, video.VideoID)

	case Playlist:
		videoId, err := GetLatestVideo(target)
		if err != nil {
			log.Error("c.Download: ", err)
			return fmt.Errorf("c.Download: %s", err)
		}
		video := Video{VideoID: videoId, DownloadMode: target.DownloadMode}
		Download(video, downloadQuality, fileExtension)
		err = UpdateLatestDownloaded(target, video.VideoID)
		if err != nil {
			log.Error("p.Download: ", err)
			return fmt.Errorf("p.Download: %s", err)
		}
		return UpdateDownloadHistory(target, video.VideoID)
		return nil
	case Video:
		log.Info("downloading video")
		var cmd *exec.Cmd
		if target.DownloadMode == "Audio Only" {
			log.Info("downloading audio only")
			if downloadQuality == "best" {
				downloadQuality = "0"
			} else if downloadQuality == "medium" {
				downloadQuality = "5"
			} else if downloadQuality == "worst" {
				downloadQuality = "9"
			}
			log.Info("download quality set to: ", downloadQuality)
			cmd = exec.Command("youtube-dl", "--extract-audio", "--audio-format", fileExtension, "--audio-quality", downloadQuality, "-o", "downloads/%(uploader)s/audio/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+target.VideoID)
		} else if target.DownloadMode == "Video And Audio" {
			cmd = exec.Command("youtube-dl", "-f", downloadQuality, "-o", "downloads/%(uploader)s/video/%(title)s.%(ext)s", "https://www.youtube.com/watch?v="+target.VideoID)
		}
		log.Info("executing youtube-dl command: ", cmd.String())
		err := cmd.Run()
		if err != nil {
			log.Error(err, cmd.String())
			return fmt.Errorf("v.Download: %s", err)
		}
		return nil
	}
	return nil
}