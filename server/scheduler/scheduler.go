package scheduler

import (
	"golty/repository"
	"golty/ytdl"

	"go.uber.org/zap"
)

type Scheduler struct {
	ytdl                 *ytdl.Ytdl
	logger               *zap.Logger
	channels             []*repository.Channel
	currentlyDownloading []*repository.Channel
}

func New(ytdl *ytdl.Ytdl, logger *zap.Logger) *Scheduler {
	channels := []*repository.Channel{}
	currentlyDownloading := []*repository.Channel{}

	return &Scheduler{ytdl: ytdl, logger: logger, channels: channels, currentlyDownloading: currentlyDownloading}
}
