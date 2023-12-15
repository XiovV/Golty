package scheduler

import (
	"fmt"
	"golty/repository"
)

type ChannelDownloadOptions struct {
	Video                    bool
	Audio                    bool
	Format                   string
	Resolution               string
	AutomaticallyDownloadNew bool
	DownloadEntire           bool
}

func (s *Scheduler) registerChannel(channel *repository.Channel) {
	s.channels = append(s.channels, channel)
}

func (s *Scheduler) DownloadChannel(channel *repository.Channel, options ChannelDownloadOptions) {
	fmt.Println("downloading channel")
}
