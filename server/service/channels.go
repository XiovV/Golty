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
	log := s.logger.With(zap.String("channelUrl", channel.ChannelUrl))
	log.Info("getting channel videos")
	channelVideos, err := s.ytdl.GetChannelVideos(channel.ChannelUrl)
	if err != nil {
		log.Error("could not get channel videos", zap.Error(err))
		return
	}

	log.Info("got channel videos successfully", zap.Int("numberOfVideos", len(channelVideos)))

	videoDownloadOptions := ytdl.VideoDownloadOptions{Video: options.Audio, Audio: options.Audio, Resolution: options.Resolution, Output: ytdl.CHANNELS_DEFAULT_OUTPUT}

	for _, videoId := range channelVideos {
		log := log.With(zap.String("videoId", videoId))

		log.Info("downloading video", zap.String("channelUrl", channel.ChannelUrl))
		err = s.ytdl.DownloadVideo(videoId, videoDownloadOptions)
		if err != nil {
			log.Error("could not download video", zap.Error(err))
			return
		}

		log.Info("video downloaded successfully")

		log.Info("getting video metadata")
		metadata, err := s.ytdl.GetVideoMetadata(videoId)
		if err != nil {
			log.Error("could not extract video metadata", zap.Error(err))
			return
		}

		log.Info("getting video filesize")
		channelPath := ytdl.CHANNELS_DIR + channel.ChannelName
		videoSize, err := s.ytdl.GetVideoSize(channelPath, videoId)
		if err != nil {
			log.Error("could not get video file size", zap.Error(err))
			return
		}

		log.Info("storing video metadata")

		video := repository.Video{
			ChannelId:    channel.ID,
			VideoId:      metadata.ID,
			Title:        metadata.Title,
			ThumbnailUrl: metadata.ThumbnailURL,
			Size:         videoSize,
		}

		err = s.repository.InsertVideo(video)
		if err != nil {
			log.Error("could not insert video metadata", zap.Error(err))
			return
		}
	}

	log.Info("channel downloaded successfully")
}
