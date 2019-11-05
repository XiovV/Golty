package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

func UpdateChannelsDatabase(channelURL string) bool {
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
		} else if channelURL == "" {
			fmt.Println("channelURL can't be empty", channelURL)
			exists = true
			break
		} else {
			exists = false
		}
	}

	if exists == true {
		fmt.Println("channel already added", channelURL)
		return true
	} else {
		fmt.Println("adding to db: ", channelURL)
		db = append(db, Channel{ChannelURL: channelURL})
		writeDb(db)
	}
	return false
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
