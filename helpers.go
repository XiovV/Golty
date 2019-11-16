package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetChannels returns the contents of channels.json
func GetChannels() ([]Channel, error) {
	log.Info("getting all channels from channels.json")
	jsonFile, err := os.Open(CONFIG_ROOT + "channels.json")
	if err != nil {
		log.Error("From GetChannels()", err)
		return []Channel{}, fmt.Errorf("From GetChannels(): %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var db []Channel

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		log.Error("From GetChannels()", err)
		return []Channel{}, fmt.Errorf("From GetChannels(): %v", err)
	}
	log.Info("successfully read all channels")
	return db, nil
}

func GetPlaylists() ([]Playlist, error) {
	log.Info("getting all playlists from playlists.json")
	jsonFile, err := os.Open(CONFIG_ROOT + "playlists.json")
	if err != nil {
		log.Error("From GetPlaylists: ", err)
		return []Playlist{}, fmt.Errorf("From GetPlaylists(): %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var playlists []Playlist

	err = json.Unmarshal(byteValue, &playlists)
	if err != nil {
		log.Error("From GetPlaylists: ", err)
		return []Playlist{}, fmt.Errorf("From GetPlaylists: %v", err)
	}
	log.Info("successfully read all playlists")
	return playlists, nil
}

// CheckAll goes through channels.json and checks for new videos
func CheckAllChannels() (Response, error) {
	log.Info("checking for all channels")
	allChannelsInDb, err := GetChannels()
	if err != nil {
		return Response{}, fmt.Errorf("From CheckAll(): %v", err)
	}
	var foundFor []string
	var preferredExtension string

	for _, item := range allChannelsInDb {
		channel := Channel{ChannelURL: item.ChannelURL}
		channel, err = channel.GetFromDatabase()
		if err != nil {
			return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the channel from database" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
		}

		if item.ChannelURL == channel.ChannelURL {
			videoId, err := channel.GetLatestVideo()
			if err != nil {
				log.Error("There was an error getting latest video: ", err)
				return Response{Type: "Error", Key: "GETTING_LATEST_VIDEO_ERROR", Message: "There was an error getting the latestvideo" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
			}

			item.UpdateLastChecked()
			if item.LatestDownloaded == videoId.VideoID {
				log.Info("no new videos found for: ", item.ChannelURL)
			} else {
				log.Info("new video detected for: ", item.ChannelURL)
				foundFor = append(foundFor, item.ChannelURL)
				if channel.DownloadMode == "Audio Only" {
					preferredExtension = channel.PreferredExtensionForAudio
				} else if channel.DownloadMode == "Video And Audio" {
					preferredExtension = channel.PreferredExtensionForVideo
				}
				go channel.Download(channel.DownloadMode, preferredExtension, "best")
				channel.UpdateLatestDownloaded(videoId.VideoID)
			}
		}
	}
	if len(foundFor) == 0 {
		return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos found."}, nil
	}
	return Response{Type: "Success", Key: "NEW_VIDEOS_FOR_CHANNELS", Message: strings.Join(foundFor, ",")}, nil
}

func CheckAllPlaylists() (Response, error) {
	log.Info("checking for all playlists")
	allPlaylistsInDb, err := GetPlaylists()
	if err != nil {
		return Response{}, fmt.Errorf("From CheckAllPlaylists(): %v", err)
	}
	var foundFor []string
	var preferredExtension string

	for _, playlist := range allPlaylistsInDb {
		p := Playlist{PlaylistURL: playlist.PlaylistURL}
		p, err = p.GetFromDatabase()
		if err != nil {
			return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the playlist from database" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
		}

		if playlist.PlaylistURL == p.PlaylistURL {
			videoId, err := p.GetLatestVideo()
			if err != nil {
				log.Error("There was an error getting latest video: ", err)
				return Response{Type: "Error", Key: "GETTING_LATEST_VIDEO_ERROR", Message: "There was an error getting the latestvideo" + err.Error()}, fmt.Errorf("CheckAll: %s", err)
			}

			playlist.UpdateLastChecked()
			if playlist.LatestDownloaded == videoId.VideoID {
				log.Info("no new videos found for: ", playlist.PlaylistURL)
			} else {
				log.Info("new video detected for: ", playlist.PlaylistURL)
				foundFor = append(foundFor, playlist.PlaylistURL)
				if p.DownloadMode == "Audio Only" {
					preferredExtension = p.PreferredExtensionForAudio
				} else if p.DownloadMode == "Video And Audio" {
					preferredExtension = p.PreferredExtensionForVideo
				}
				go p.Download(p.DownloadMode, preferredExtension, "best")
				p.UpdateLatestDownloaded(videoId.VideoID)
			}
		}
	}
	if len(foundFor) == 0 {
		return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos found."}, nil
	}
	return Response{Type: "Success", Key: "NEW_VIDEOS_FOR_PLAYLISTS", Message: strings.Join(foundFor, ",")}, nil
}

// CheckNow requires c.ChannelURL
func (c Channel) CheckNow() (Response, error) {
	log.Info("checking for new videos")
	allChannelsInDb, err := GetChannels()
	if err != nil {
		return Response{}, fmt.Errorf("From CheckNow(): %v", err)
	}

	var preferredExtension string

	channel, err := c.GetFromDatabase()
	if err != nil {
		return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the channel from database" + err.Error()}, fmt.Errorf("CheckNow: %s", err)
	}
	channelURL := c.ChannelURL

	channelMetadata, err := channel.GetMetadata()
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting channel metadata: " + err.Error()}, nil
	}

	err = c.UpdateLastChecked()
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_UPDATING_LAST_CHECKED", Message: "There was an error updating latest checked date and time: " + err.Error()}, nil
	}

	for _, item := range allChannelsInDb {
		if item.ChannelURL == channelURL {
			if item.LatestDownloaded == channelMetadata.ID {
				log.Info("no new videos found for: ", channelURL)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected for " + item.Name}, nil
			} else {
				log.Info("new video detected for: ", channelURL)
				if channel.DownloadMode == "Audio Only" {
					preferredExtension = channel.PreferredExtensionForAudio
				} else if channel.DownloadMode == "Video And Audio" {
					preferredExtension = channel.PreferredExtensionForVideo
				}
				err := channel.Download(channel.DownloadMode, preferredExtension, "best")
				if err != nil {
					log.Error(err)
					return Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: err.Error()}, nil
				}
				channel.UpdateLatestDownloaded(channelMetadata.ID)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected for " + item.Name}, nil
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}, nil
}

func (p Playlist) CheckNow() (Response, error) {
	log.Info("checking for new videos")
	allPlaylistsInDb, err := GetPlaylists()
	if err != nil {
		return Response{}, fmt.Errorf("From p.CheckNow(): %v", err)
	}

	var preferredExtension string

	playlist, err := p.GetFromDatabase()
	if err != nil {
		return Response{Type: "Error", Key: "GETTING_FROM_DATABASE_ERROR", Message: "There was an error getting the playlist from database" + err.Error()}, fmt.Errorf("CheckNow: %s", err)
	}
	playlistURL := p.PlaylistURL

	playlistMetadata, err := playlist.GetMetadata()
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_GETTING_METADATA", Message: "There was an error getting playlist metadata: " + err.Error()}, nil
	}

	err = p.UpdateLastChecked()
	if err != nil {
		return Response{Type: "Error", Key: "ERROR_UPDATING_LAST_CHECKED", Message: "There was an error updating latest checked date and time: " + err.Error()}, nil
	}

	for _, playlist := range allPlaylistsInDb {
		if playlist.PlaylistURL == playlistURL {
			if playlist.LatestDownloaded == playlistMetadata.ID {
				log.Info("no new videos found for: ", playlistURL)
				return Response{Type: "Success", Key: "NO_NEW_VIDEOS", Message: "No new videos detected for " + playlist.Name}, nil
			} else {
				log.Info("new video detected for: ", playlistURL)
				if playlist.DownloadMode == "Audio Only" {
					preferredExtension = playlist.PreferredExtensionForAudio
				} else if playlist.DownloadMode == "Video And Audio" {
					preferredExtension = playlist.PreferredExtensionForVideo
				}
				err := playlist.Download(playlist.DownloadMode, preferredExtension, "best")
				if err != nil {
					log.Error(err)
					return Response{Type: "Error", Key: "ERROR_DOWNLOADING_VIDEO", Message: err.Error()}, nil
				}
				playlist.UpdateLatestDownloaded(playlistMetadata.ID)
				return Response{Type: "Success", Key: "NEW_VIDEO_DETECTED", Message: "New video detected for " + playlist.Name + " and downloaded"}, nil
			}
		}
	}
	log.Error("Something went terribly wrong")
	return Response{Type: "Error", Key: "UNKNOWN_ERROR", Message: "Something went wrong"}, nil
}

func CreateDirIfNotExist(dirName string) {
	log.Info("creating channel directory")
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			log.Error("Couldn't create channel directory: ", err)
		} else {
			log.Info("Channel directory created successfully")
		}
	}
}

func RemoveAtIndexChannel(s []Channel, index int) []Channel {
	return append(s[:index], s[index+1:]...)
}

func RemoveAtIndexPlaylist(s []Playlist, index int) []Playlist {
	return append(s[:index], s[index+1:]...)
}

func GetChannelName(channelURL string) string {
	return strings.Split(channelURL, "/")[4]
}

func ReturnResponse(w http.ResponseWriter, res Response) {
	log.Info("returning response: ", res)
	json.NewEncoder(w).Encode(res)
}
