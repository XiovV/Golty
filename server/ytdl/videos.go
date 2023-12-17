package ytdl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type VideoMetadata struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	ThumbnailURL   string `json:"thumbnail"`
	UploadDate     string `json:"upload_date"`
	DurationString string `json:"duration_string"`
}

func (y *Ytdl) DownloadVideo(videoId string, options VideoDownloadOptions) error {
	if options.Video {
		resolutionFlag := fmt.Sprintf("res:%s", y.resolutions[options.Resolution])
		_, err := y.exec("-S", resolutionFlag, "-o", options.Output, videoId)
		return err
	}

	return nil
}

func (y *Ytdl) GetVideoMetadata(videoId string) (VideoMetadata, error) {
	out, err := y.exec("--print", "%(.{title,thumbnail,id,upload_date,duration_string})#+j", videoId)
	if err != nil {
		return VideoMetadata{}, err
	}

	var videoMetadata VideoMetadata
	err = json.Unmarshal(out, &videoMetadata)
	if err != nil {
		return VideoMetadata{}, err
	}

	return videoMetadata, nil
}

func (y *Ytdl) getVideoFilename(path, videoId string) (string, error) {
	videos, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, video := range videos {
		if strings.HasPrefix(video.Name(), videoId) {
			return video.Name(), nil
		}
	}

	return "", errors.New("video not found")
}

func (y *Ytdl) GetVideoSize(path, videoId string) (int64, error) {
	videoFilename, err := y.getVideoFilename(path, videoId)
	if err != nil {
		return 0, err
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", path, videoFilename))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	size := fileStat.Size()

	return size, nil
}
