package service

import (
	"fmt"
	"golty/repository"
	"golty/ytdl"
	"time"

	"go.uber.org/zap"
)

type ChannelsService struct {
	repository           *repository.Repository
	logger               *zap.Logger
	ytdl                 *ytdl.Ytdl
	channels             []*repository.Channel
	currentlyDownloading []*repository.Channel
}

type ChannelState struct {
	*repository.Channel
	State string `json:"state"`
}

type ChannelDownloadOptions struct {
	Video              bool
	Audio              bool
	Format             string
	Resolution         string
	DownloadNewUploads bool
	DownloadEntire     bool
}

func NewChannelsService(repository *repository.Repository, logger *zap.Logger, ytdl *ytdl.Ytdl) *ChannelsService {
	return &ChannelsService{repository: repository, logger: logger, ytdl: ytdl}
}

func (s *ChannelsService) DownloadChannel(channel repository.Channel, options ChannelDownloadOptions) {
	log := s.logger.With(zap.String("channelUrl", channel.ChannelUrl))

	s.addToCurrentlyDownloading(&channel)
	defer s.removeFromCurrentlyDownloading(&channel)

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

		s.addToCurrentlyDownloading(&channel)

		videoDownloadOptions := ytdl.VideoDownloadOptions{
			Video:      s.integerToBoolean(channelSettings.DownloadVideo),
			Audio:      s.integerToBoolean(channelSettings.DownloadAudio),
			Resolution: channelSettings.Resolution,
			Format:     channelSettings.Format,
			Output:     ytdl.CHANNELS_DEFAULT_OUTPUT,
		}

		s.downloadChannelVideos(channel, missingVideos, videoDownloadOptions)

		s.removeFromCurrentlyDownloading(&channel)
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
			ThumbnailUrl: metadata.ThumbnailURL,
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

func (s *ChannelsService) addToCurrentlyDownloading(channel *repository.Channel) {
	s.logger.Debug("adding channel to currently downloading", zap.String("channelHandle", channel.ChannelHandle))
	s.currentlyDownloading = append(s.currentlyDownloading, channel)
}

func (s *ChannelsService) removeFromCurrentlyDownloading(channel *repository.Channel) {
	s.logger.Debug("removing channel from currently downloading", zap.String("channelHandle", channel.ChannelHandle))
	for i, c := range s.currentlyDownloading {
		if c.ChannelHandle == channel.ChannelHandle {
			s.currentlyDownloading = append(s.currentlyDownloading[:i], s.currentlyDownloading[i+1:]...)
		}
	}
}

func (s *ChannelsService) isChannelCurrentlyDownloading(channelId int) bool {
	for _, channel := range s.currentlyDownloading {
		fmt.Println(channel.ID)
		if channel.ID == channelId {
			return true
		}
	}

	return false
}

func (s *ChannelsService) GetChannelDownloadState(channelId int) string {
	if s.isChannelCurrentlyDownloading(channelId) {
		return "downloading"
	}

	return "idle"
}
