package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/go-get-youtube/youtube"
)

func DownloadVideoAndAudio(videoID, videoTitle string) {
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
	video.Download(0, videoTitle+".mp4", option)
}

func DownloadAudio(videoID, videoTitle, channelName string) {
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
	video.Download(0, videoTitle+".mp4", option)
	fmt.Println("Removing mp4...")
	os.Remove(videoTitle + ".mp4")
}

func DownloadVideo(videoID, videoTitle, channelName string) {
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
	video.Download(0, channelName+"/"+videoTitle+".mp4", option)
}

func Download(channelName, channelType, downloadMode string) {
	if channelType == "user" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle, channelName)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle, channelName)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle, channelName)
		}
	} else if channelType == "channel" {
		videoId, videoTitle := GetLatestVideo(channelName, channelType)
		if downloadMode == "Audio Only" {
			DownloadAudio(videoId, videoTitle, channelName)
		} else if downloadMode == "Video Only" {
			DownloadVideo(videoId, videoTitle, channelName)
		} else if downloadMode == "Video And Audio" {
			DownloadVideo(videoId, videoTitle, channelName)
		}
	}
}
