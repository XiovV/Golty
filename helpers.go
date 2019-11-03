package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func UpdateLatestDownloaded(channelName, videoID string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	jsonFile, err := os.Open("channels.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	json.Unmarshal(byteValue, &db)

	for i, item := range db {
		if strings.Contains(item.ChannelURL, channelName) {
			db[i].LatestDownloaded = videoID
			break
		}
	}

	writeDb(db)
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
			exists = true
			break
		} else {
			exists = false
		}
	}

	if exists == true {
		fmt.Println("channel already added")
	} else {
		db = append(db, Channel{ChannelURL: channelURL})

		writeDb(db)
	}
}

func writeDb(db []Channel) {
	result, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(result, &db)

	file, _ := json.MarshalIndent(db, "", " ")

	_ = ioutil.WriteFile("channels.json", file, 0644)
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

func UploadChecker() {
	for {
		time.Sleep(10 * time.Second)
		go CheckNow(nil, "")
	}
}

func CheckNow(channels []string, channelType string) {
	allChannelsInDb := GetChannels()

	if channels == nil {
		for _, item := range allChannelsInDb {
			channelName := strings.Split(item.ChannelURL, "/")[4]
			channelType := strings.Split(item.ChannelURL, "/")[3]

			fmt.Println(channelName, channelType)

			if strings.Contains(item.ChannelURL, channelName) {
				videoId, videoTitle := GetLatestVideo(channelName, channelType)

				if item.LatestDownloaded == videoId {
					fmt.Println("NOW NEW VIDEOS DETECTED FOR: ", item.ChannelURL)
				} else {
					fmt.Println("DOWNLOAD FOR: ", item.ChannelURL)
					DownloadAudio(videoId, videoTitle)
					UpdateLatestDownloaded(item.ChannelURL, videoId)
				}
			}

		}
	} else {
		videoId, videoTitle := GetLatestVideo(channels[0], channelType)

		for _, item := range allChannelsInDb {
			if strings.Contains(item.ChannelURL, channels[0]) {
				if item.LatestDownloaded == videoId {
					break
				} else {
					DownloadAudio(videoId, videoTitle)
					UpdateLatestDownloaded(channels[0], videoId)
				}
			}
		}
	}
}

func ReturnResponse(w http.ResponseWriter, response string) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	channels := GetChannels()
	t.Execute(w, Response{Channels: channels, Status: response})
}
