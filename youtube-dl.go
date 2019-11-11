package main

import (
	"encoding/json"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func getUserMetadata(channelName string) Video {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", "https://www.youtube.com/user/"+channelName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	metaData := &ChannelInformation{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Fatal(err)
	}

	return Video{VideoID: metaData.ID}
}

func getChannelMetadata(channelName string) Video {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", "https://www.youtube.com/channel/"+channelName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	metaData := &ChannelInformation{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Fatal(err)
	}

	return Video{VideoID: metaData.ID}
}

func GetLatestVideo(channelName, channelType string) Video {
	if channelType == "user" {
		return getUserMetadata(channelName)
	}
	return getChannelMetadata(channelName)
}

func (v Video) DownloadYTDL() error {
	cmd := exec.Command("youtube-dl", "-o", "downloads/ %(uploader)s/ %(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
		return err
	}

	return nil
}

func (v Video) DownloadAudioYTDL() error {
	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "-o", "downloads/ %(uploader)s/ %(title)s.%(ext)s", "https://www.youtube.com/watch?v="+v.VideoID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
		return err
	}

	return nil
}
