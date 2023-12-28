package service

import (
	"golty/repository"
	"io"
	"net/http"
	"os"
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
