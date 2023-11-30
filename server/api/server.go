package api

import (
	"golty/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Logger *zap.Logger
}

func New(config *config.Config, logger *zap.Logger) *Server {
	return &Server{Config: config, Logger: logger}
}

func (s *Server) Start() error {
	e := echo.New()
	e.Use(middleware.CORS())

	v1 := e.Group("/v1")
	usersPublic := v1.Group("/users")

	usersPublic.POST("/login", s.loginUserHandler)

	return e.Start(":" + s.Config.Port)
}
