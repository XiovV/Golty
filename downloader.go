package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/go-get-youtube/youtube"
)

func DownloadAudio(videoID, videoTitle, channelName string) {
	path := channelName + "/" + videoTitle + ".mp4"

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	video, err := youtube.Get(videoID)
	if err != nil {
		log.Panic(err)
	}

	option := &youtube.Option{
		Rename: false,
		Resume: true,
		Mp3:    true,
	}
	video.Download(0, path, option)
	fmt.Println("Removing mp4...")
	os.Remove(path)
}

func DownloadVideo(videoID, videoTitle, channelName string) {
	path := channelName + "/" + videoTitle + ".mp4"
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	video, err := youtube.Get(videoID)
	if err != nil {
		log.Panic(err)
	}

	option := &youtube.Option{
		Rename: false,
		Resume: true,
		Mp3:    false,
	}
	video.Download(0, path, option)
}

func Download(channelName, channelType, downloadMode string) {
	if channelType == "user" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle, channelName)
			UpdateLatestDownloaded(channelName, videoId)
		} else if downloadMode == "Video Only" {
			UpdateLatestDownloaded(channelName, videoId)
			DownloadVideo(videoId, videoTitle, channelName)
		} else if downloadMode == "Video And Audio" {
			UpdateLatestDownloaded(channelName, videoId)
			DownloadVideo(videoId, videoTitle, channelName)
		}
	} else if channelType == "channel" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			UpdateLatestDownloaded(channelName, videoId)
			DownloadAudio(videoId, videoTitle, channelName)
		} else if downloadMode == "Video Only" {
			UpdateLatestDownloaded(channelName, videoId)
			DownloadVideo(videoId, videoTitle, channelName)
		} else if downloadMode == "Video And Audio" {
			UpdateLatestDownloaded(channelName, videoId)
			DownloadVideo(videoId, videoTitle, channelName)
		}
	}
}
