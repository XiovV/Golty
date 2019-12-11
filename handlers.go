package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type version struct {
	VersionNumber string
}

func HandleGetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version{VersionNumber: VERSION})
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "log.log")
}

func HandlePlaylists(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/playlists.html")
}

func HandleVideos(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/videos.html")
}

func HandleAddTarget(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to add a target")
	w.Header().Set("Content-Type", "application/json")
	var targetData AddTargetPayload
	err := json.NewDecoder(r.Body).Decode(&targetData)
	if err != nil {
		errRes = Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
	}
	//sanitiseInputURL(targetData.URL)
	target := DownloadTarget{URL: targetData.URL, Type: targetData.Type}

	doesTargetExist, err := target.DoesExist()
	if err != nil {
		log.Info("error doesChannelExist: ", err)
		errRes = Response{Type: "Error", Key: "DOES_EXIST_ERROR", Message: "There was an error while trying to check if the channel already exists" + err.Error()}
	}

	if doesTargetExist {
		log.Info("this playlist already exists")
		okRes = Response{Type: "Success", Key: "PLAYLIST_ALREADY_EXISTS", Message: "This playlists already exists"}
	} else {
		targetMetadata, err := target.GetMetadata()
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting channel metadata: " + err.Error()}
		}

		if targetData.DownloadMode == "Audio Only" {
			target = DownloadTarget{URL: targetData.URL, DownloadMode: targetData.DownloadMode, Name: targetMetadata.Playlist, PreferredExtensionForAudio: targetData.FileExtension, DownloadHistory: []string{}, LastChecked: time.Now().Format("01-02-2006 15:04:05"), CheckingInterval: "", Type: targetData.Type, DownloadPath: targetData.DownloadPath}
		} else if targetData.DownloadMode == "Video And Audio" {
			target = DownloadTarget{URL: targetData.URL, DownloadMode: targetData.DownloadMode, Name: targetMetadata.Playlist, PreferredExtensionForVideo: targetData.FileExtension, DownloadHistory: []string{}, LastChecked: time.Now().Format("01-02-2006 15:04:05"), CheckingInterval: "", Type: targetData.Type, DownloadPath: targetData.DownloadPath}
		}

		err = target.AddToDatabase()
		if err != nil {
			log.Error(err)
			errRes = Response{Type: "Error", Key: "ERROR_ADDING_PLAYLIST", Message: "There was an error adding the playlist to the database" + err.Error()}
		}
		err = target.Download(targetData.DownloadQuality, targetData.FileExtension, targetData.DownloadEntire)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_DOWNLOADING", Message: "There was an error downloading the target " + err.Error()}
		}
		err = target.UpdateLatestDownloaded(targetMetadata.ID)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_UPDATING_LATEST_DOWNLOADED", Message: "There was an error while updating the playlist's latest downloaded video id: " + err.Error()}
		}
		err = target.UpdateDownloadHistory(targetMetadata.ID)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_ERROR_UPDATING_DOWNLOAD_HISTORY", Message: "There was an error while updating the playlist's download history: " + err.Error()}
		}
		okRes = Response{Type: "Success", Key: "ADD_PLAYLIST_SUCCESS", Message: "Successfully added and downloaded latest video"}
	}
	if errRes.Type == "Error" {
		ReturnResponse(w, errRes)
	} else if okRes.Type == "Success" {
		ReturnResponse(w, okRes)
	}
}

func HandleCheckTarget(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check a playlist for new uploads")
	w.Header().Set("Content-Type", "application/json")
	var targetData AddTargetPayload
	err := json.NewDecoder(r.Body).Decode(&targetData)
	if err != nil {
		errRes = Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
		ReturnResponse(w, errRes)
		return
	}
	target := DownloadTarget{URL: targetData.URL, Type: targetData.Type}
	target, _ = target.GetFromDatabase()

	newVideoFound, videoId, err := target.CheckNow()
	if err != nil {
		errRes = Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLIST", Message: "There was an error while checking the playlist: " + err.Error()}
		ReturnResponse(w, errRes)
		return
	}
	if newVideoFound == true {
		err = target.Download("best", targetData.FileExtension, false)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLIST", Message: "There was an error while checking the channel: " + err.Error()}
			ReturnResponse(w, errRes)
			return
		}
		err = target.UpdateLatestDownloaded(videoId)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLIST", Message: "There was an error while checking the channel: " + err.Error()}
			ReturnResponse(w, errRes)
			return
		}
		err = target.UpdateDownloadHistory(videoId)
		if err != nil {
			errRes = Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLIST", Message: "There was an error while checking the channel: " + err.Error()}
		}
		okRes = Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected for " + target.Name + " and downloaded"}
		ReturnResponse(w, okRes)
		return
	} else {
		okRes = Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected for " + target.Name}
		ReturnResponse(w, okRes)
		return
	}
}

func HandleDeleteTarget(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to delete a playlist")

	w.Header().Set("Content-Type", "application/json")
	var targetData DeleteTargetPayload
	err := json.NewDecoder(r.Body).Decode(&targetData)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
		return
	}
	targetURL := targetData.URL
	targetURL = strings.Replace(targetURL, "delTarget", "", -1)
	target := DownloadTarget{URL: targetURL, Type: targetData.Type}

	target.Delete()

	ReturnResponse(w, Response{Type: "Success", Key: "DELETE_PLAYLIST_SUCCESS", Message: "Successfully removed"})
	return
}

func HandleGetTargets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var targetData GetTargetPayload
	err := json.NewDecoder(r.Body).Decode(&targetData)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
		return
	}
	log.Info("got this data: ", targetData)
	targets, err := GetAll(targetData.Type)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting playlists: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(targets)
	return
}

func HandleCheckAllTargets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var targetData GetTargetPayload
	err := json.NewDecoder(r.Body).Decode(&targetData)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
		return
	}
	res, err := CheckAll(targetData.Type)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLISTS", Message: "There was an error while checking playlists: " + err.Error()})
		return
	}
	ReturnResponse(w, res)
	return
}

func HandleUpdateCheckingInterval(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to update the checking interval")
	w.Header().Set("Content-Type", "application/json")

	var interval CheckingIntervalPayload
	err := json.NewDecoder(r.Body).Decode(&interval)
	if err != nil {
		errRes = Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()}
		ReturnResponse(w, errRes)
		return
	}
	target := DownloadTarget{Type: interval.Type}
	res, err := target.UpdateCheckingInterval(interval.CheckingInterval)
	if err != nil {
		errRes = Response{Type: "Error", Key: "ERROR_UPDATING_CHECKING_INTERVAL", Message: "There was an updating the checking interval: " + err.Error()}
		ReturnResponse(w, errRes)
		return
	}
	ReturnResponse(w, res)
}
