package ytdl

import (
	"encoding/json"
	"fmt"
)

type VideoMetadata struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail"`
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
	out, err := y.exec("--print", "%(.{title,thumbnail,id})#+j", videoId)
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
