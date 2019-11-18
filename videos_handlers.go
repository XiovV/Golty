package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleDownloadVideo(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to download a video")
	var ytdlCommand string
	var videoData DownloadVideoPayload
	err := json.NewDecoder(r.Body).Decode(&videoData)
	if err != nil {
		log.Error("HandleDownloadVideo: ", err)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	if videoData.DownloadMode == "Audio Only" {
		ytdlCommand = "youtube-dl -f bestaudio[ext="+videoData.FileExtension+"] -o downloads/videos/%(uploader)s/audio/%(title)s.%(ext)s "+videoData.VideoURL
	} else if videoData.DownloadMode == "Video And Audio" {
		ytdlCommand = "youtube-dl -o downloads/videos/%(uploader)s/video/%(title)s.%(ext)s "+videoData.VideoURL
	}
	DownloadVideo(ytdlCommand)

	videoData.AddToDatabase()

	ReturnResponse(w, Response{Type: "Success", Key: "DOWNLOAD_VIDEO_SUCCESS", Message: "Video successfully downloaded."})

}

func HandleGetVideos(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to get all videos")
	w.Header().Set("Content-Type", "application/json")

	videos, err := GetVideos()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_GETTING_VIDEOS", Message: "There was an error while getting videos: " + err.Error()}
		json.NewEncoder(w).Encode(res)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting channels: " + err.Error()})
	}
	json.NewEncoder(w).Encode(videos)
}
