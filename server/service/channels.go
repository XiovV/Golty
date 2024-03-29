package service

import (
	"fmt"
	"golty/queue"
	"golty/repository"
	"golty/ytdl"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ChannelsService struct {
	repository                  *repository.Repository
	logger                      *zap.Logger
	ytdl                        *ytdl.Ytdl
	channels                    []*repository.Channel
	currentlyDownloading        *repository.Channel
	CurrentlyDownloadingChannel chan *repository.Channel
	ChannelsQueue               *queue.ChannelsQueue
}

type ChannelDownloadOptions struct {
	Video          bool
	Audio          bool
	Format         string
	Quality        string
	DownloadEntire bool
	Sync           bool
}

func NewChannelsService(repo *repository.Repository, logger *zap.Logger, ytdl *ytdl.Ytdl, channelsQueue *queue.ChannelsQueue) *ChannelsService {
	currentlyDownloadingChannel := make(chan *repository.Channel)

	return &ChannelsService{
		repository:                  repo,
		logger:                      logger,
		ytdl:                        ytdl,
		CurrentlyDownloadingChannel: currentlyDownloadingChannel,
		ChannelsQueue:               channelsQueue,
	}
}

func (s *ChannelsService) DownloadChannel(channel repository.Channel, options ChannelDownloadOptions) {
	log := s.logger.With(zap.String("channelUrl", channel.ChannelUrl))

	log.Debug("getting channel videos")

	channelVideos, err := s.ytdl.GetChannelVideos(channel.ChannelUrl)
	if err != nil {
		log.Error("could not get channel videos", zap.Error(err))
		return
	}

	videoDownloadOptions := s.channelOptionsToVideoOptions(options, ytdl.CHANNELS_DEFAULT_OUTPUT)
	err = s.DownloadChannelVideos(channel, channelVideos, videoDownloadOptions)
	if err != nil {
		log.Error("could not download channel", zap.Error(err))
		return
	}

	log.Info("channel downloaded successfully")
}

func (s *ChannelsService) AddChannel(channel repository.Channel, downloadOptions ChannelDownloadOptions) error {
	log := s.logger.With(zap.String("channelUrl", channel.ChannelUrl))

	avatarDestination := fmt.Sprintf("avatars/%s.png", channel.ChannelHandle)
	err := s.downloadImage(channel.AvatarUrl, avatarDestination)
	if err != nil {
		log.Error("could not download image", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	channel.AvatarUrl = avatarDestination

	log.Debug("inserting channel into the database")
	createdChannel, err := s.repository.InsertChannel(channel)
	if err != nil {
		log.Error("could not insert channel", zap.Error(err))

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return echo.NewHTTPError(http.StatusBadRequest, "This channel already exists!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Debug("persisting channel download settings")

	err = s.repository.InsertChannelDownloadSettings(s.channelOptionsToDBChannelOptions(createdChannel.ID, downloadOptions))
	if err != nil {
		log.Error("could not persist channel download settings", zap.Error(err))
		return err
	}

	if !downloadOptions.DownloadEntire {
		return nil
	}

	log.Debug("enqueueing channel")
	s.ChannelsQueue.Enqueue(&createdChannel)

	return nil
}

func (s *ChannelsService) ResumeDownloads() {
	channelDownloadSettings, err := s.repository.GetAllDownloadSettings()
	if err != nil {
		s.logger.Error("unable to get channel download settings", zap.Error(err))
		return
	}

	for _, channelSettings := range channelDownloadSettings {
		log := s.logger.With(zap.Int("channelId", channelSettings.ChannelId))

		// TODO: use a JOIN here
		log.Debug("resuming downloads")
		channel, err := s.repository.FindChannelByID(channelSettings.ChannelId)
		if err != nil {
			log.Error("could not find channel by id", zap.Error(err))
			return
		}

		log.Debug("getting missing videos")
		missingVideos, err := s.GetMissingVideos(channel)
		if err != nil {
			log.Error("could not get missing videos", zap.Error(err))
			return
		}

		if len(missingVideos) == 0 && bool(channelSettings.Sync) {
			log.Debug("no missing videos, continuing")
			continue
		}

		s.ChannelsQueue.Enqueue(&channel)
	}
}

func (s *ChannelsService) DownloadChannelVideos(channel repository.Channel, videos []string, options ytdl.VideoDownloadOptions) error {
	for _, videoId := range videos {
		log := s.logger.With(zap.String("videoId", videoId))

		log.Debug("downloading video")
		err := s.ytdl.DownloadVideo(videoId, options)
		if err != nil {
			return err
		}
		dateDownloaded := time.Now().Unix()

		log.Debug("video downloaded successfully")

		log.Debug("getting video metadata")
		metadata, err := s.ytdl.GetVideoMetadata(videoId)
		if err != nil {
			return err
		}

		thumbnailDestination := fmt.Sprintf("thumbnails/%s_thumbnail.jpg", metadata.ID)
		err = s.downloadImage(metadata.ThumbnailURL, thumbnailDestination)
		if err != nil {
			return nil
		}

		parseUploadDate, err := time.Parse("20060102", metadata.UploadDate)
		if err != nil {
			return err
		}

		log.Debug("getting video filesize")
		channelPath := ytdl.CHANNELS_DIR + channel.ChannelName
		videoSize, err := s.ytdl.GetVideoSize(channelPath, videoId)
		if err != nil {
			return err
		}

		log.Debug("storing video metadata")

		video := repository.Video{
			ChannelId:    channel.ID,
			VideoId:      metadata.ID,
			Title:        metadata.Title,
			ThumbnailUrl: thumbnailDestination,
			Size:         videoSize,
			DownloadDate: dateDownloaded,
			UploadDate:   parseUploadDate.Unix(),
			Duration:     metadata.DurationString,
		}

		err = s.repository.InsertVideo(video)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ChannelsService) SyncChannel(channelId int) (int, error) {
	channel, err := s.repository.FindChannelByID(channelId)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "could not find channel")
	}

	missingVideos, err := s.GetMissingVideos(channel)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go s.ChannelsQueue.Enqueue(&channel)

	return len(missingVideos), nil
}
