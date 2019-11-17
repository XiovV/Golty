package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func HandleAddPlaylist(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to add a playlist")
	w.Header().Set("Content-Type", "application/json")
	var playlistData AddPlaylistPayload
	err := json.NewDecoder(r.Body).Decode(&playlistData)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}

	playlist := Playlist{PlaylistURL: playlistData.PlaylistURL}

	doesPlaylistExist, err := DoesExist(playlist)
	if err != nil {
		log.Info("error doesChannelExist: ", err)
		ReturnResponse(w, Response{Type: "Error", Key: "DOES_EXIST_ERROR", Message: "There was an error while trying to check if the channel already exists" + err.Error()})
	}

	if doesPlaylistExist == true {
		log.Info("this playlist already exists")
		ReturnResponse(w, Response{Type: "Success", Key: "PLAYLIST_ALREADY_EXISTS", Message: "This playlists already exists"})
	} else {
		playlistMetadata, err := GetMetadata(playlist)
		if err != nil {
			ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting channel metadata: " + err.Error()})
		}

		switch playlistMetadata := playlistMetadata.(type) {
		case PlaylistMetadata:
			if playlistData.DownloadMode == "Audio Only" {
				playlist = Playlist{PlaylistURL: playlistData.PlaylistURL, DownloadMode: playlistData.DownloadMode, Name: playlistMetadata.Playlist, PreferredExtensionForAudio: playlistData.FileExtension, DownloadHistory: []string{}, LastChecked: time.Now().Format("01-02-2006 15:04:05"), CheckingInterval: ""}
			} else if playlistData.DownloadMode == "Video And Audio" {
				playlist = Playlist{PlaylistURL: playlistData.PlaylistURL, DownloadMode: playlistData.DownloadMode, Name: playlistMetadata.Playlist, PreferredExtensionForVideo: playlistData.FileExtension, DownloadHistory: []string{}, LastChecked: time.Now().Format("01-02-2006 15:04:05"), CheckingInterval: ""}
			}
		}

		err = AddToDatabase(playlist)
		if err != nil {
			log.Error(err)
			ReturnResponse(w, Response{Type: "Error", Key: "ERROR_ADDING_PLAYLIST", Message: "There was an error adding the playlist to the database" + err.Error()})
		}
		if playlistData.DownloadEntireChannel == true {
			//err := playlist.DownloadEntire()
			//if err != nil {
			//	ReturnResponse(w, Response{Type: "Error", Key: "ERROR_DOWNLOADING_ENTIRE_PLAYLIST", Message: "There was an error downloading the entire playlist" + err.Error()})
			//}
		} else {
			if err != nil {
				log.Error(err)
				ReturnResponse(w, Response{Type: "Error", Key: "ERROR_ADDING_PLAYLIST", Message: "There was an error adding the playlist to the database" + err.Error()})
			}
			err = Download(playlist, playlistData.DownloadQuality, playlistData.FileExtension)
			if err != nil {
				log.Error(err)
				ReturnResponse(w, Response{Type: "Error", Key: "ERROR_DOWNLOADING", Message: "There was an error while downloading: " + err.Error()})
			}

			ReturnResponse(w, Response{Type: "Success", Key: "ADD_PLAYLIST_SUCCESS", Message: "Playlist successfully added and downloaded latest video"})
		}
	}
}

func HandleCheckPlaylist(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check a playlist for new uploads")
	w.Header().Set("Content-Type", "application/json")
	var data AddPlaylistPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	playlist := Playlist{PlaylistURL: data.PlaylistURL}

	res, err := playlist.CheckNow()
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLIST", Message: "There was an error while checking the playlist: " + err.Error()})
	}
	ReturnResponse(w, res)
}

func HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to get all playlists")
	w.Header().Set("Content-Type", "application/json")

	playlists, err := GetPlaylists()
	if err != nil {
		res := Response{Type: "Error", Key: "ERROR_GETTING_PLAYLISTS", Message: "There was an error while getting playlists: " + err.Error()}
		json.NewEncoder(w).Encode(res)
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_GETTING_CHANNELS", Message: "There was an error while getting playlists: " + err.Error()})
	}
	json.NewEncoder(w).Encode(playlists)
}

func HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to delete a playlist")

	w.Header().Set("Content-Type", "application/json")
	var data DeletePlaylistPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_PARSING_DATA", Message: "There was an error parsing json: " + err.Error()})
	}
	playlistURL := data.PlaylistURL
	playlistURL = strings.Replace(playlistURL, "delPlaylist", "", -1)
	playlist := Playlist{PlaylistURL: playlistURL}

	Delete(playlist)

	ReturnResponse(w, Response{Type: "Success", Key: "DELETE_PLAYLIST_SUCCESS", Message: "Playlist removed"})
}

func HandleCheckAllPlaylists(w http.ResponseWriter, r *http.Request) {
	log.Info("received a request to check all playlists for new uploads")
	w.Header().Set("Content-Type", "application/json")
	res, err := CheckAllPlaylists()
	if err != nil {
		ReturnResponse(w, Response{Type: "Error", Key: "ERROR_CHECKING_PLAYLISTS", Message: "There was an error while checking playlists: " + err.Error()})
	}
	ReturnResponse(w, res)
}
