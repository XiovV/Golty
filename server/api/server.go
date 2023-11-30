package api

import (
	"golty/config"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Logger *zap.Logger
}

func New(config *config.Config, logger *zap.Logger) *Server {
	return &Server{Config: config, Logger: logger}
}

func (s *Server) Router() *echo.Echo {
	e := echo.New()

	v1 := e.Group("/v1")
	usersPublic := v1.Group("/users")

	usersPublic.POST("/login", s.loginUserHandler)

	return e
}
