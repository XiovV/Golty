package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/addchannel", HandleAddChannel).Methods("POST")

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))

	// channelURL := "https://www.youtube.com/user/NewRetroWave"

	// channelName := strings.Split(channelURL, "/")[4]
	// channelType := strings.Split(channelURL, "/")[3]

	// if channelType == "user" {
	// 	uploadsId := GetUserUploadsID(channelName)
	// 	videoId, videoTitle := GetUserVideoData(uploadsId)

	// 	DownloadVideo(videoId, videoTitle)
	// } else if channelType == "channel" {
	// 	videoId, videoTitle := GetChannelVideoData(channelName)
	// 	fmt.Println(videoId, videoTitle)
	// }

}
