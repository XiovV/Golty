package service

import (
	"golty/ytdl"

	"go.uber.org/zap"
)

func (s *ChannelsService) StartQueueConsumer() {
	s.logger.Debug("started queue consumer")
	for {
		channel := s.ChannelsQueue.GetFirst()
		if channel == nil {
			s.logger.Fatal("received nil from the queue. This is a fatal error and must be reported")
		}

		s.logger.Debug("got channel in queue", zap.Int("channelId", channel.ID))

		channelSettings, err := s.repository.GetChannelDownloadSettings(channel.ID)
		if err != nil {
			s.logger.Error("could not get channel download settings", zap.Error(err), zap.Int("channelId", channel.ID))
			continue
		}

		missingVideos, err := s.GetMissingVideos(*channel)
		if err != nil {
			s.logger.Error("could not get missing videos", zap.Error(err))
			continue
		}

		videoDownloadOptions := ytdl.VideoDownloadOptions{
			Video:   bool(channelSettings.Video),
			Audio:   bool(channelSettings.Audio),
			Quality: channelSettings.Quality,
			Format:  channelSettings.Format,
			Output:  ytdl.CHANNELS_DEFAULT_OUTPUT,
		}

		s.logger.Debug("downloading channel", zap.Int("channelId", channel.ID))
		s.DownloadChannelVideos(*channel, missingVideos, videoDownloadOptions)

		s.logger.Debug("channel finished downloading", zap.Int("channelId", channel.ID))
		s.ChannelsQueue.Dequeue()
	}
}
