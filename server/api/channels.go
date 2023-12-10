package api

import (
	"golty/repository"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) getChannelInfo(c echo.Context) error {
	channelName := c.Param("channelUrl")

	channelInfo, err := s.Ytdl.GetChannelInfo(channelName)
	if err != nil {
		s.Logger.Error("unable to get channel info", zap.Error(err), zap.String("channelName", channelName))
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

	err := s.Repository.InsertChannel(channel)
	if err != nil {
		s.Logger.Error("could not insert channel", zap.Error(err), zap.String("channelName", addChannelRequest.Channel.ChannelName))

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return echo.NewHTTPError(http.StatusBadRequest, "This channel already exists!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)
}
