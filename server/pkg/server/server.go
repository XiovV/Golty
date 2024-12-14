package server

import (
	"go.uber.org/zap"

	"github.com/XiovV/Golty/pkg/config"
	"github.com/labstack/echo/v4"
)

func New(config *config.Config, logger *zap.SugaredLogger) error {
	e := echo.New()
	e.HideBanner = true

	return e.Start(":" + config.Port)
}
