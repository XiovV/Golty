package api

import (
	"golty/config"
	"golty/repository"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Config     *config.Config
	Logger     *zap.Logger
	Repository *repository.Repository
}

func New(config *config.Config, logger *zap.Logger, repository *repository.Repository) *Server {
	return &Server{Config: config, Logger: logger, Repository: repository}
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS(), middleware.Logger())

	v1 := e.Group("/v1")
	usersPublic := v1.Group("/users")

	usersPublic.Use(echojwt.WithConfig(jwtConfig))
	usersPublic.POST("/login", s.loginUserHandler)

	return e.Start(":" + s.Config.Port)
}
