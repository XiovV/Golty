package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleDownloadVideo(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to download a video")
	var videoData DownloadVideoPayload
	err := json.NewDecoder(r.Body).Decode(&videoData)
	if err != nil {
		log.Error("HandleDownloadVideo: ", err)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}

	err = videoData.Download()
	if err != nil {
		log.Error("HandleDownloadVideo: ", err)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: "There was an error downloading the video: " + err.Error()})
	}

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
