package service

import (
	"golty/repository"
	"golty/ytdl"
	"time"

	"go.uber.org/zap"
)

type ChannelsService struct {
	repository    *repository.Repository
	logger        *zap.Logger
	ytdl          *ytdl.Ytdl
	channels      []*repository.Channel
	channelStates []ChannelState
}

type ChannelState struct {
	*repository.Channel
	State string `json:"state"`
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

	log.Info("getting missing videos")
	missingVideos, err := s.getMissingVideos(channel)
	if err != nil {
		log.Error("could not get missing videos", zap.Error(err))
		return
	}

	if len(missingVideos) == 0 {
		log.Info("no new videos detected")
		return
	}

	log.Info("got missing videos successfully", zap.Int("numberOfVideos", len(missingVideos)))

	videoDownloadOptions := ytdl.VideoDownloadOptions{Video: options.Audio, Audio: options.Audio, Resolution: options.Resolution, Output: ytdl.CHANNELS_DEFAULT_OUTPUT}

	for _, videoId := range missingVideos {
		log := log.With(zap.String("videoId", videoId))

		log.Info("downloading video", zap.String("channelUrl", channel.ChannelUrl))
		err = s.ytdl.DownloadVideo(videoId, videoDownloadOptions)
		if err != nil {
			log.Error("could not download video", zap.Error(err))
			return
		}
		dateDownloaded := time.Now().Unix()

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
			ChannelId:      channel.ID,
			VideoId:        metadata.ID,
			Title:          metadata.Title,
			ThumbnailUrl:   metadata.ThumbnailURL,
			Size:           videoSize,
			DateDownloaded: dateDownloaded,
		}

		err = s.repository.InsertVideo(video)
		if err != nil {
			log.Error("could not insert video metadata", zap.Error(err))
			return
		}
	}

	log.Info("channel downloaded successfully")
}

func (s *ChannelsService) getMissingVideos(channel repository.Channel) ([]string, error) {
	downloadedVideos, err := s.repository.GetChannelVideos(channel.ID)
	if err != nil {
		return []string{}, err
	}

	channelVideos, err := s.ytdl.GetChannelVideos(channel.ChannelUrl)
	if err != nil {
		return []string{}, err
	}

	if len(downloadedVideos) == len(channelVideos) {
		return []string{}, nil
	}

	numberOfMissingVideos := len(channelVideos) - len(downloadedVideos)

	return channelVideos[:numberOfMissingVideos], nil
}
