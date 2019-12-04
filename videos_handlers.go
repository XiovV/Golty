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
		errRes = Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
	}
	log.Info(videoData)
	if videoData.DownloadPath == "" {
		if videoData.DownloadMode == "Audio Only" {
			ytdlCommand = "youtube-dl -f bestaudio[ext=" + videoData.FileExtension + "] -o " + videoData.DownloadPath + " " + videoData.VideoURL
		} else if videoData.DownloadMode == "Video And Audio" {
			ytdlCommand = "youtube-dl -o " + videoData.DownloadPath + " " + videoData.VideoURL
		}
	} else {
		if videoData.DownloadMode == "Audio Only" {
			ytdlCommand = "youtube-dl -f bestaudio[ext=" + videoData.FileExtension + "] -o downloads" + videoData.DownloadPath + " " + videoData.VideoURL
		} else if videoData.DownloadMode == "Video And Audio" {
			ytdlCommand = "youtube-dl -o downloads" + videoData.DownloadPath + " " + videoData.VideoURL
		}
	}

	err = DownloadVideo(ytdlCommand)
	if err != nil {
		errRes = Response{Type: "Error", Key: "DOWNLOAD_VIDEO_ERROR", Message: "There was an error while downloading the video: " + err.Error()}
	}
	err = videoData.AddToDatabase()
	if err != nil {
		errRes = Response{Type: "Error", Key: "ADDING_TO_DATABASE_ERROR", Message: "There was an error while adding the video into the database: " + err.Error()}
	}
	if errRes.Type == "Error" {
		ReturnResponse(w, errRes)
	} else if errRes.Type == "" {
		ReturnResponse(w, Response{Type: "Success", Key: "DOWNLOAD_VIDEO_SUCCESS", Message: "Video successfully downloaded."})
	}
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
