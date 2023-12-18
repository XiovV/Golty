package api

import (
	"golty/config"
	"golty/repository"
	"golty/service"
	"golty/ytdl"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Config          *config.Config
	Logger          *zap.Logger
	Repository      *repository.Repository
	ChannelsService *service.ChannelsService
	Ytdl            *ytdl.Ytdl
}

func New(config *config.Config, logger *zap.Logger, repository *repository.Repository, channelsService *service.ChannelsService, ytdl *ytdl.Ytdl) *Server {
	return &Server{Config: config, Logger: logger, Repository: repository, ChannelsService: channelsService, Ytdl: ytdl}
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS(), middleware.Logger())

	v1 := e.Group("/v1")
	usersPublic := v1.Group("/users")
	{
		usersPublic.POST("/login", s.loginUserHandler)
	}

	usersAuth := v1.Group("/users")
	usersAuth.Use(echojwt.WithConfig(jwtConfig))
	{
		usersAuth.GET("/me", s.getLoggedInUser)
	}

	channels := v1.Group("/channels")
	// channels.Use(echojwt.WithConfig(jwtConfig))
	{
		channels.GET("", s.getChannelsHandler)
		channels.POST("", s.addChannelHandler)
		channels.GET("/info/:channelUrl", s.getChannelInfoHandler)
		channels.GET("/:channelHandle", s.getChannelHandler)
		channels.GET("/videos/:channelHandle", s.getChannelVideosHandler)
		channels.GET("/state", s.getChannelStateWs)
	}

	return e.Start(":" + s.Config.Port)
}
