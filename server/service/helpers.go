package service

import (
	"golty/repository"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

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

func (s *ChannelsService) GetMissingVideos(channel repository.Channel) ([]string, error) {
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

func (s *ChannelsService) registerChannel(channel *repository.Channel) {
	s.channels = append(s.channels, channel)
}
