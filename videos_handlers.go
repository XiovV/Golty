package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func HandleDownloadVideo(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to download a video")
	var ytdlCommand YTDLCommand
	var videoData DownloadVideoPayload
	err := json.NewDecoder(r.Body).Decode(&videoData)
	if err != nil {
		log.Error("HandleDownloadVideo: ", err)
		errRes = Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
	}
	log.Info(videoData)
	ytdlCommand = YTDLCommand{
		Binary: "youtube-dl",
		Target: videoData.VideoURL,
	}
	switch videoData.DownloadPath {
	default:
		ytdlCommand.Output = filepath.Join(dlRoot, videoData.DownloadPath)
	case "":
		ytdlCommand.Output = filepath.Join(dlRoot, "/%(uploader)s/audio/%(title)s.%(ext)s")
	}
	switch videoData.DownloadMode {
	case "Audio Only":
		ytdlCommand.FileType = "bestaudio[ext=" + videoData.FileExtension + "]"
	case "Video And Audio":
		ytdlCommand.FileType = "bestvideo[height<=" + videoData.DownloadQuality + "]" + "+ bestaudio/best[height<=" + videoData.DownloadQuality + "]"
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
