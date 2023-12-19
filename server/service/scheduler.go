package service

import (
	"golty/ytdl"
	"time"

	"go.uber.org/zap"
)

func (s *ChannelsService) StartScheduler() {
	ticker := time.Tick(30 * time.Second)
	for {
		select {
		case <-ticker:
			s.checkChannels()
		}
	}
}

func (s *ChannelsService) checkChannels() {
	for _, channel := range s.channels {
		missingVideos, err := s.getMissingVideos(*channel)
		if err != nil {
			s.logger.Error("could not get missing vidoes", zap.Error(err))
			continue
		}

		if len(missingVideos) == 0 {
			continue
		}

		channelSettings, err := s.repository.GetChannelDownloadSettings(channel.ID)

		downloadOptions := ytdl.VideoDownloadOptions{
			Video:      s.integerToBoolean(channelSettings.DownloadVideo),
			Audio:      s.integerToBoolean(channelSettings.DownloadAudio),
			Resolution: channelSettings.Resolution,
			Format:     channelSettings.Format,
			Output:     ytdl.CHANNELS_DEFAULT_OUTPUT,
		}

		s.downloadChannelVideos(*channel, missingVideos, downloadOptions)
	}
}
