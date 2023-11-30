package main

import (
	"golty/api"
	"golty/config"
	zapLogger "golty/logger"
	"log"
	"net/http"

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

	server := api.New(c, logger)
	router := server.Router()

	logger.Info("server is listening...", zap.String("port", c.Port))

	err = http.ListenAndServe(":"+c.Port, router)
	if err != nil {
		logger.Error("failed to initialise server", zap.Error(err))
	}
}
