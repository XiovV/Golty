package ytdl

import "fmt"

func (y *Ytdl) DownloadVideo(videoId string, options VideoDownloadOptions) error {
	if options.Video {
		resolutionFlag := fmt.Sprintf("res:%s", y.resolutions[options.Resolution])
		_, err := y.exec("-S", resolutionFlag, "-o", options.Output, videoId)
		return err
	}

	return nil
}
