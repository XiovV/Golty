package server

import (
	"go.uber.org/zap"

	"github.com/XiovV/Golty/pkg/config"
	"github.com/XiovV/Golty/pkg/db"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config    *config.Config
	logger    *zap.SugaredLogger
	db        *db.DB
	jwtConfig echojwt.Config
}

func New(config *config.Config, logger *zap.SugaredLogger, db *db.DB) *Server {
	return &Server{config: config, logger: logger, db: db}
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true

	api := e.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.POST("/login", s.loginHandler)
	authGroup.POST("/refresh-token", s.refreshTokenHandler)

	// protected := api.Group("/protected")
	// protected.Use(echojwt.WithConfig(s.jwtConfig))
	// protected.GET("/test", s.protectedRouteHandler)

	return e.Start(":" + s.config.Port)
}
