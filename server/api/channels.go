package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) getChannelInfo(c echo.Context) error {
	channelName := c.Param("channelName")

	channelInfo, err := s.Ytdl.GetChannelInfo(channelName)
	if err != nil {
		s.Logger.Error("unable to get channel info", zap.Error(err), zap.String("channelName", channelName))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, channelInfo)
}
