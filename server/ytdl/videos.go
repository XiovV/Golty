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
	optionsFlag := y.generateDownloadOptionsFlag(options)

	flags := []string{}

	for _, flag := range optionsFlag {
		flags = append(flags, flag)
	}

	flags = append(flags, []string{"-o", options.Output, videoId}...)

	_, err := y.exec(flags...)
	if err != nil {
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

func (y *Ytdl) generateDownloadOptionsFlag(options VideoDownloadOptions) []string {
	// download video and audio
	if options.Video && options.Audio {
		return []string{"-f", fmt.Sprintf("bv*[height<=%s]+ba", options.Quality)}
	}

	// download video only
	if options.Video && !options.Audio {
		return []string{"-f", fmt.Sprintf("bv[height<=%s]", options.Quality)}
	}

	// download audio only
	if !options.Video && options.Audio {
		format := options.Format
		if format == "auto" {
			format = "mp3"
		}
		return []string{"-x", "--audio-format", format, "--audio-quality", options.Quality}
	}

	// returning an empty string because I feel like returning nil here could be dangerous
	return []string{""}
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
