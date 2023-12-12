package ytdl

import (
	"go.uber.org/zap"
)

const CHANNELS_DEFAULT_OUTPUT = "videos/channels/%(uploader)s/%(id)s.%(ext)s"

type ChannelInfo struct {
	UploaderID string `json:"uploader_id"`
	Uploader   string `json:"uploader"`
	Avatar     struct {
		URL string `json:"url"`
	} `json:"avatar"`
}

type ChannelDownloadOptions struct {
	Video                    bool
	Audio                    bool
	Format                   string
	Resolution               string
	AutomaticallyDownloadNew bool
	DownloadEntire           bool
}

type ChannelVideos struct {
	Videos []string `json:"videos"`
}

func (y *Ytdl) GetChannelInfo(channelUrl string) (ChannelInfo, error) {
	ytdlOutput, err := y.exec("--playlist-items", "0", "-J", channelUrl)
	if err != nil {
		return ChannelInfo{}, err
	}

	var channelInfo ChannelInfo
	err = y.jq(string(ytdlOutput), &channelInfo, "{uploader_id, uploader, avatar: (.thumbnails[] | select(.id==\"avatar_uncropped\"))}")
	if err != nil {
		return ChannelInfo{}, err
	}

	return channelInfo, err
}

func (y *Ytdl) GetChannelVideos(channelUrl string) ([]string, error) {
	ytdlOutput, err := y.exec("--get-id", "--flat-playlist", channelUrl)
	if err != nil {
		return nil, err
	}

	var channelVideos ChannelVideos
	err = y.jq(string(ytdlOutput), &channelVideos, "-nR", "{ \"videos\": [inputs] }")
	if err != nil {
		return nil, err
	}

	return channelVideos.Videos, nil
}

func (y *Ytdl) DownloadChannel(channelUrl string, downloadOptions ChannelDownloadOptions) {
	y.logger.Info("getting channel videos", zap.String("channelUrl", channelUrl))
	channelVideos, err := y.GetChannelVideos(channelUrl)
	if err != nil {
		y.logger.Error("could not get channel videos", zap.Error(err), zap.String("channelUrl", channelUrl))
		return
	}

	y.logger.Info("got channel videos successfully", zap.String("channelUrl", channelUrl), zap.Int("numberOfVideos", len(channelVideos)))

	videoDownloadOptions := VideoDownloadOptions{Video: downloadOptions.Video, Audio: downloadOptions.Audio, Resolution: downloadOptions.Resolution, Output: CHANNELS_DEFAULT_OUTPUT}

	for _, video := range channelVideos {
		y.logger.Info("downloading video", zap.String("channelUrl", channelUrl), zap.String("videoId", video))
		err = y.downloadVideo(video, videoDownloadOptions)
		if err != nil {
			y.logger.Error("downloading video failed", zap.Error(err), zap.String("channelUrl", channelUrl), zap.String("videoId", video))
			return
		}

		y.logger.Info("video downloaded successfully", zap.String("channelUrl", channelUrl), zap.String("videoId", video))
	}
}
