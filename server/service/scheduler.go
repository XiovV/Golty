package service

import (
	"golty/ytdl"
	"time"

	"go.uber.org/zap"
)

const SCHEDULER_DEFAULT_INTERVAL = 2 * time.Hour

func (s *ChannelsService) StartScheduler() {
	ticker := time.Tick(SCHEDULER_DEFAULT_INTERVAL)
	for {
		select {
		case <-ticker:
			s.checkChannels()
		}
	}
}

// TODO: this needs a complete rewrite
// Add a function that will return only the channels that have the downloadNewUploads setting enabled
// Go through each channel and enqueue it if it's got more than 1 missing video
func (s *ChannelsService) checkChannels() {
	for _, channel := range s.channels {
		missingVideos, err := s.GetMissingVideos(*channel)
		if err != nil {
			s.logger.Error("could not get missing vidoes", zap.Error(err))
			continue
		}

		if len(missingVideos) == 0 {
			continue
		}

		channelSettings, err := s.repository.GetChannelDownloadSettings(channel.ID)

		downloadOptions := ytdl.VideoDownloadOptions{
			Video:      bool(channelSettings.DownloadVideo),
			Audio:      bool(channelSettings.DownloadAudio),
			Resolution: channelSettings.Resolution,
			Format:     channelSettings.Format,
			Output:     ytdl.CHANNELS_DEFAULT_OUTPUT,
		}

		s.DownloadChannelVideos(*channel, missingVideos, downloadOptions)
	}
}
