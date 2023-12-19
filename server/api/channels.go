package api

import (
	"golty/repository"
	"golty/service"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) getChannelInfoHandler(c echo.Context) error {
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

func (s *Server) addChannelHandler(c echo.Context) error {
	var addChannelRequest struct {
		Channel struct {
			ChannelUrl    string `json:"channelUrl"`
			ChannelName   string `json:"channelName"`
			ChannelHandle string `json:"channelHandle"`
			AvatarUrl     string `json:"avatarUrl"`
		} `json:"channel"`
		DownloadSettings struct {
			DownloadVideo      bool   `json:"downloadVideo"`
			DownloadAudio      bool   `json:"downloadAudio"`
			Resolution         string `json:"resolution"`
			Format             string `json:"format"`
			DownloadNewUploads bool   `json:"downloadNewUploads"`
			DownloadEntire     bool   `json:"downloadEntire"`
		} `json:"downloadSettings"`
	}

	if err := c.Bind(&addChannelRequest); err != nil {
		s.Logger.Error("json input is invalid", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	if !addChannelRequest.DownloadSettings.DownloadVideo && !addChannelRequest.DownloadSettings.DownloadAudio {
		return echo.NewHTTPError(http.StatusBadRequest, "Either the Video or the Audio checkbox must be checked, or both.")
	}

	if !addChannelRequest.DownloadSettings.DownloadNewUploads && !addChannelRequest.DownloadSettings.DownloadEntire {
		return echo.NewHTTPError(http.StatusBadRequest, "Please toggle either the 'Automatically download new uploads' or 'Download the entire channel' switch, or both.")
	}

	channel := repository.Channel{
		ChannelUrl:    addChannelRequest.Channel.ChannelUrl,
		ChannelName:   addChannelRequest.Channel.ChannelName,
		ChannelHandle: addChannelRequest.Channel.ChannelHandle,
		AvatarUrl:     addChannelRequest.Channel.AvatarUrl,
	}

	channelDownloadOptions := service.ChannelDownloadOptions{
		Video:              addChannelRequest.DownloadSettings.DownloadVideo,
		Audio:              addChannelRequest.DownloadSettings.DownloadAudio,
		Resolution:         addChannelRequest.DownloadSettings.Resolution,
		Format:             addChannelRequest.DownloadSettings.Format,
		DownloadEntire:     addChannelRequest.DownloadSettings.DownloadEntire,
		DownloadNewUploads: addChannelRequest.DownloadSettings.DownloadNewUploads,
	}

	err := s.ChannelsService.AddChannel(channel, channelDownloadOptions)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

var isDownloading bool

func (s *Server) getChannelsHandler(c echo.Context) error {
	channels, err := s.Repository.GetChannelsWithSize()
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
		TotalSize     int    `json:"totalSize"`
		State         string `json:"state"`
	}

	channelResponses := []channelResponse{}
	for _, channel := range channels {
		channelResponses = append(channelResponses, channelResponse{
			ID:            channel.ID,
			ChannelName:   channel.ChannelName,
			ChannelHandle: channel.ChannelHandle,
			ChannelUrl:    channel.ChannelUrl,
			AvatarUrl:     channel.AvatarUrl,
			TotalVideos:   channel.TotalVideos,
			TotalSize:     channel.TotalSize,
			State:         "idle",
		})
	}

	return c.JSON(http.StatusOK, channelResponses)
}

func (s *Server) getChannelHandler(c echo.Context) error {
	channelHandle := strings.Replace(c.Param("channelHandle"), "%40", "@", 1)

	channel, err := s.Repository.FindChannelByHandleWithSize(channelHandle)
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
		TotalSize     int    `json:"totalSize"`
		State         string `json:"state"`
	}

	response := channelResponse{
		ID:            channel.ID,
		ChannelName:   channel.ChannelName,
		ChannelHandle: channel.ChannelHandle,
		ChannelUrl:    channel.ChannelUrl,
		AvatarUrl:     channel.AvatarUrl,
		TotalVideos:   channel.TotalVideos,
		TotalSize:     channel.TotalSize,
		State:         "downloading",
	}

	return c.JSON(http.StatusOK, response)
}

func (s *Server) getChannelVideosHandler(c echo.Context) error {
	channelHandle := strings.Replace(c.Param("channelHandle"), "%40", "@", 1)
	log := s.Logger.With(zap.String("channelHandle", channelHandle))

	channel, err := s.Repository.FindChannelByHandleWithSize(channelHandle)
	if err != nil {
		log.Error("could not find channel", zap.Error(err))
		return echo.NewHTTPError(http.StatusNotFound, "channel does not exist ")
	}

	channelVideos, err := s.Repository.GetChannelVideos(channel.ID)
	if err != nil {
		log.Error("could not get channel videos", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	type channelVideosResponse struct {
		VideoId      string `json:"videoId"`
		Title        string `json:"title"`
		ThumbnailUrl string `json:"thumbnailUrl"`
		Size         int64  `json:"size"`
		DownloadDate int64  `json:"downloadDate"`
		Duration     string `json:"duration"`
	}

	response := []channelVideosResponse{}
	for _, video := range channelVideos {
		response = append(response, channelVideosResponse{
			VideoId:      video.VideoId,
			Title:        video.Title,
			ThumbnailUrl: video.ThumbnailUrl,
			Size:         video.Size,
			DownloadDate: video.DownloadDate,
			Duration:     video.Duration,
		})
	}

	return c.JSON(http.StatusOK, response)
}

var upgrader = websocket.Upgrader{}

func (s *Server) getChannelStateWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	type channelStateResponse struct {
		ChannelId int    `json:"channelId"`
		State     string `json:"state"`
	}

	var channelId int
	go func() {
		for {
			channel, ok := <-s.ChannelsService.CurrentlyDownloadingChannel
			if !ok {
				return
			}

			if channel == nil {

				err := ws.WriteJSON(channelStateResponse{ChannelId: channelId, State: "idle"})
				if err != nil {
					c.Logger().Error(err)
				}
				continue
			}

			channelId = channel.ID

			err := ws.WriteJSON(channelStateResponse{ChannelId: channelId, State: "downloading"})
			if err != nil {
				c.Logger().Error(err)
			}

		}
	}()

	return nil
}
