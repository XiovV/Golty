package server

import (
	"go.uber.org/zap"

	"github.com/XiovV/Golty/pkg/config"
	"github.com/XiovV/Golty/pkg/db"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config *config.Config
	logger *zap.SugaredLogger
	db     *db.DB
}

func New(config *config.Config, logger *zap.SugaredLogger, db *db.DB) *Server {
	return &Server{config: config, logger: logger, db: db}
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true

	return e.Start(":" + s.config.Port)
}
