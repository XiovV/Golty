package main

import (
	"golty/api"
	"golty/config"
	zapLogger "golty/logger"
	"golty/repository"
	"golty/ytdl"
	"log"

	"go.uber.org/zap"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatalln("config err:", err)
	}

	logger, err := zapLogger.Init(c.Environment)
	if err != nil {
		log.Fatalln("could not init logger:", err)
	}

	repository, err := repository.New(c.SQLiteDir)
	if err != nil {
		logger.Fatal("could not init database", zap.Error(err))
	}

	err = initServer(repository, logger)
	if err != nil {
		logger.Fatal("could not initialise server", zap.Error(err))
	}

	ytdl := ytdl.New("yt-dlp", logger)

	server := api.New(c, logger, repository, ytdl)

	err = server.Start()
	if err != nil {
		logger.Error("failed to initialise server", zap.Error(err))
	}
}
