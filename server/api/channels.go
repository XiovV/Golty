package api

import (
	"golty/repository"
	"golty/ytdl"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) getChannelInfo(c echo.Context) error {
	channelUrl := c.Param("channelUrl")

	channelUrlSplit := strings.Split(channelUrl, "/")
	channelName := channelUrlSplit[len(channelUrlSplit)-1]
	if strings.Contains(channelName, "@") && len(channelName) == 1 || len(channelName) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "channel name or handle must be provided")
	}

	channelInfo, err := s.Ytdl.GetChannelInfo(channelUrl)
	if err != nil {
		s.Logger.Error("unable to get channel info", zap.Error(err), zap.String("channelUrl", channelUrl))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, channelInfo)
}

func (s *Server) addChannel(c echo.Context) error {
	var addChannelRequest struct {
		Channel struct {
			ChannelUrl    string `json:"channelUrl"`
			ChannelName   string `json:"channelName"`
			ChannelHandle string `json:"channelHandle"`
			AvatarUrl     string `json:"avatarUrl"`
		} `json:"channel"`
		DownloadSettings struct {
			DownloadVideo         bool   `json:"downloadVideo"`
			DownloadAudio         bool   `json:"downloadAudio"`
			Resolution            string `json:"resolution"`
			Format                string `json:"format"`
			DownloadAutomatically bool   `json:"downloadAutomatically"`
			DownloadEntireChannel bool   `json:"downloadEntireChannel"`
		} `json:"downloadSettings"`
	}

	if err := c.Bind(&addChannelRequest); err != nil {
		s.Logger.Error("json input is invalid", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	channel := repository.Channel{
		ChannelUrl:    addChannelRequest.Channel.ChannelUrl,
		ChannelName:   addChannelRequest.Channel.ChannelName,
		ChannelHandle: addChannelRequest.Channel.ChannelHandle,
		AvatarUrl:     addChannelRequest.Channel.AvatarUrl,
	}

	// err := s.Repository.InsertChannel(channel)
	// if err != nil {
	// 	s.Logger.Error("could not insert channel", zap.Error(err), zap.String("channelName", addChannelRequest.Channel.ChannelName))
	//
	// 	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
	// 		return echo.NewHTTPError(http.StatusBadRequest, "This channel already exists!")
	// 	}
	//
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	s.Logger.Info("downloading channel", zap.String("channelName", channel.ChannelName))
	go s.Ytdl.DownloadChannel(channel.ChannelUrl, ytdl.ChannelDownloadOptions{})

	return c.NoContent(http.StatusCreated)
}

func (s *Server) getChannels(c echo.Context) error {
	channels, err := s.Repository.GetChannels()
	if err != nil {
		s.Logger.Error("could not get all channels", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	type channelResponse struct {
		ID            int    `json:"id"`
		ChannelName   string `json:"channelName"`
		ChannelHandle string `json:"channelHandle"`
		ChannelUrl    string `json:"channelUrl"`
		AvatarUrl     string `json:"avatarUrl"`
		TotalVideos   int    `json:"totalVideos"`
		TotalSize     string `json:"totalSize"`
	}

	channelResponses := []channelResponse{}
	for _, channel := range channels {
		channelResponses = append(channelResponses, channelResponse{
			ID:            channel.ID,
			ChannelName:   channel.ChannelName,
			ChannelHandle: channel.ChannelHandle,
			ChannelUrl:    channel.ChannelUrl,
			AvatarUrl:     channel.AvatarUrl,
			TotalVideos:   0,
			TotalSize:     "0 GB",
		})
	}

	return c.JSON(http.StatusOK, channelResponses)
}

func (s *Server) getChannel(c echo.Context) error {
	channelHandle := strings.Replace(c.Param("channelHandle"), "%40", "@", 1)

	channel, err := s.Repository.GetChannelByHandle(channelHandle)
	if err != nil {
		s.Logger.Error("could not get channel by handle", zap.Error(err), zap.String("channelHandle", channelHandle))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	type channelResponse struct {
		ID            int    `json:"id"`
		ChannelName   string `json:"channelName"`
		ChannelHandle string `json:"channelHandle"`
		ChannelUrl    string `json:"channelUrl"`
		AvatarUrl     string `json:"avatarUrl"`
		TotalVideos   int    `json:"totalVideos"`
		TotalSize     string `json:"totalSize"`
	}

	return c.JSON(http.StatusOK, channelResponse{ID: channel.ID, ChannelName: channel.ChannelName, ChannelHandle: channel.ChannelHandle, ChannelUrl: channel.ChannelUrl, AvatarUrl: channel.AvatarUrl, TotalVideos: 0, TotalSize: "0 GB"})
}
