package service

import (
	"fmt"
	"golty/repository"
	"golty/ytdl"
	"io"
	"net/http"
	"os"
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
}

type ChannelDownloadOptions struct {
	Video              bool
	Audio              bool
	Format             string
	Resolution         string
	DownloadNewUploads bool
	DownloadEntire     bool
}

func NewChannelsService(repo *repository.Repository, logger *zap.Logger, ytdl *ytdl.Ytdl) *ChannelsService {
	currentlyDownloadingChannel := make(chan *repository.Channel)

	return &ChannelsService{repository: repo, logger: logger, ytdl: ytdl, CurrentlyDownloadingChannel: currentlyDownloadingChannel}
}

func (s *ChannelsService) DownloadChannel(channel repository.Channel, options ChannelDownloadOptions) {
	log := s.logger.With(zap.String("channelUrl", channel.ChannelUrl))

	s.setCurrentlyDownloading(&channel)
	defer s.unsetCurrentlyDownloading()

	log.Debug("persisting channel download settings")

	err := s.repository.InsertChannelDownloadSettings(repository.ChannelDownloadSettings{
		ChannelId:          channel.ID,
		Resolution:         options.Resolution,
		Format:             options.Format,
		DownloadVideo:      s.booleanToInteger(options.Video),
		DownloadAudio:      s.booleanToInteger(options.Audio),
		DownloadEntire:     s.booleanToInteger(options.DownloadEntire),
		DownloadNewUploads: s.booleanToInteger(options.DownloadNewUploads),
	})
	if err != nil {
		log.Error("could not persist channel download settings", zap.Error(err))
		return
	}

	log.Debug("getting channel videos")

	channelVideos, err := s.ytdl.GetChannelVideos(channel.ChannelUrl)
	if err != nil {
		log.Error("could not get channel videos", zap.Error(err))
		return
	}

	videoDownloadOptions := ytdl.VideoDownloadOptions{Video: options.Audio, Audio: options.Audio, Resolution: options.Resolution, Output: ytdl.CHANNELS_DEFAULT_OUTPUT}

	err = s.downloadChannelVideos(channel, channelVideos, videoDownloadOptions)
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
	log.Debug("channel successfully inserted into the database")

	log.Info("downloading channel")

	go s.DownloadChannel(createdChannel, downloadOptions)

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

		log.Debug("resuming downloads")
		channel, err := s.repository.FindChannelByID(channelSettings.ChannelId)
		if err != nil {
			log.Error("could not find channel by id", zap.Error(err))
			return
		}

		log.Debug("getting missing videos")
		missingVideos, err := s.getMissingVideos(channel)
		if err != nil {
			log.Error("could not get missing videos", zap.Error(err))
			return
		}

		if len(missingVideos) == 0 {
			log.Debug("no missing videos, continuing")
			continue
		}

		s.setCurrentlyDownloading(&channel)

		videoDownloadOptions := ytdl.VideoDownloadOptions{
			Video:      s.integerToBoolean(channelSettings.DownloadVideo),
			Audio:      s.integerToBoolean(channelSettings.DownloadAudio),
			Resolution: channelSettings.Resolution,
			Format:     channelSettings.Format,
			Output:     ytdl.CHANNELS_DEFAULT_OUTPUT,
		}

		s.downloadChannelVideos(channel, missingVideos, videoDownloadOptions)

		s.unsetCurrentlyDownloading()
	}
}

func (s *ChannelsService) downloadChannelVideos(channel repository.Channel, videos []string, options ytdl.VideoDownloadOptions) error {
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

func (s *ChannelsService) downloadImage(imageUrl string, destination string) error {
	response, err := http.Get(imageUrl)
	if err != nil {
		return nil
	}

	defer response.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
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

func (s *ChannelsService) integerToBoolean(integer int) bool {
	if integer > 0 {
		return true
	}

	return false
}

func (s *ChannelsService) booleanToInteger(boolean bool) int {
	if boolean {
		return 1
	}

	return 0
}

func (s *ChannelsService) setCurrentlyDownloading(channel *repository.Channel) {
	s.logger.Debug("adding channel to currently downloading", zap.String("channelHandle", channel.ChannelHandle))
	// s.currentlyDownloading = channel
	// s.CurrentlyDownloadingChannel <- channel
}

func (s *ChannelsService) unsetCurrentlyDownloading() {
	s.logger.Debug("removing channel from currently downloading")
	// s.currentlyDownloading = nil
	//
	// s.CurrentlyDownloadingChannel <- nil
}

func (s *ChannelsService) isChannelCurrentlyDownloading() bool {
	if s.currentlyDownloading == nil {
		return false
	}

	return true
}
