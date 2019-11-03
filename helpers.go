package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func DownloadAudio(videoID, videoTitle string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	video, err := youtube.Get(videoID)
	if err != nil {
		log.Panic("asdadasd", err)
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

func DownloadVideo(videoID, videoTitle string) {
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

func UpdateChannelsDatabase(channelURL string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	var exists bool

	for _, v := range db {
		if v.ChannelURL == channelURL {
			fmt.Println("already exists:", channelURL)
			exists = true
			break
		} else {
			fmt.Println("doesnt exist:", channelURL)
			exists = false
		}
	}

	if exists == true {
		fmt.Println("channel already added")
	} else {
		db = append(db, Channel{ChannelURL: channelURL})

		result, err := json.Marshal(db)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(string(result))

		json.Unmarshal(result, &db)

		file, _ := json.MarshalIndent(db, "", " ")

		_ = ioutil.WriteFile("channels.json", file, 0644)
	}
}

func GetChannels() []Channel {
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	return db
}
