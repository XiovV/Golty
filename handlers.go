package main

import (
	"log"
	"net/http"
	"text/template"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}

func HandleAddChannel(w http.ResponseWriter, r *http.Request) {
	channelURL := r.FormValue("channelURL")
	UpdateChannelsDatabase(channelURL)

	// db := ReadDatabase()

	// newDb := append(db, Database{Channels: []Channel{Channel{ChannelURL: "test2"}}})
	// fmt.Println("newDb: ", newDb)
	// downloadMode := r.FormValue("mode")

	// data := Database{
	// 	Channels: []Channel{
	// 		Channel{ChannelURL: channelURL},
	// 	},
	// 	DownloadedVideos: []DownloadedVideo{
	// 		DownloadedVideo{VideoID: "adasdasd"},
	// 	},
	// }

	// file, _ := json.MarshalIndent(data, "", " ")

	// _ = ioutil.WriteFile("database.json", file, 0644)
}
