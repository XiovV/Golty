package api

import (
	"golty/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Logger *zap.Logger
}

func New(config *config.Config, logger *zap.Logger) *Server {
	return &Server{Config: config, Logger: logger}
}

func (s *Server) Router() *gin.Engine {
	if s.Config.Environment == PRODUCTION_ENV {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), s.CORS())

	v1 := router.Group("v1")

	usersPublic := v1.Group("users")
	{
		usersPublic.POST("/login", s.loginUserHandler)
	}

	return router
}
