package service

import (
	"golty/repository"
	"golty/ytdl"

	"go.uber.org/zap"
)

type ChannelsService struct {
	repository *repository.Repository
	logger     *zap.Logger
	ytdl       *ytdl.Ytdl
}

type ChannelDownloadOptions struct {
	Video                    bool
	Audio                    bool
	Format                   string
	Resolution               string
	AutomaticallyDownloadNew bool
	DownloadEntire           bool
}

func NewChannelsService(repository *repository.Repository, logger *zap.Logger, ytdl *ytdl.Ytdl) *ChannelsService {
	return &ChannelsService{repository: repository, logger: logger, ytdl: ytdl}
}

func (s *ChannelsService) DownloadChannel(channel repository.Channel, options ChannelDownloadOptions) {
	s.logger.Info("getting channel videos", zap.String("channelUrl", channel.ChannelUrl))
	channelVideos, err := s.ytdl.GetChannelVideos(channel.ChannelUrl)
	if err != nil {
		s.logger.Error("could not get channel videos", zap.Error(err), zap.String("channelUrl", channel.ChannelUrl))
		return
	}

	s.logger.Info("got channel videos successfully", zap.String("channelUrl", channel.ChannelUrl), zap.Int("numberOfVideos", len(channelVideos)))

	videoDownloadOptions := ytdl.VideoDownloadOptions{Video: options.Audio, Audio: options.Audio, Resolution: options.Resolution, Output: ytdl.CHANNELS_DEFAULT_OUTPUT}

	for _, video := range channelVideos {
		s.logger.Info("downloading video", zap.String("channelUrl", channel.ChannelUrl), zap.String("videoId", video))
		err = s.ytdl.DownloadVideo(video, videoDownloadOptions)
		if err != nil {
			s.logger.Error("downloading video failed", zap.Error(err), zap.String("channelUrl", channel.ChannelUrl), zap.String("videoId", video))
			return
		}

		s.logger.Info("video downloaded successfully", zap.String("channelUrl", channel.ChannelUrl), zap.String("videoId", video))
	}
}
