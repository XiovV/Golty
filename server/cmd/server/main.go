package main

import (
	"log"

	"github.com/XiovV/Golty/pkg/config"
	zapLogger "github.com/XiovV/Golty/pkg/logger"
	"github.com/XiovV/Golty/pkg/server"
)

func main() {
	logger, err := zapLogger.New()
	if err != nil {
		log.Fatalln("logger init error:", err)
	}

	c, err := config.New()
	if err != nil {
		logger.Error("config error", "error", err)
	}

	err = server.New(c, logger)
	if err != nil {
		logger.Error("could not start server", "error", err)
	}
}
