package queue

import (
	"golty/repository"
	"sync"
)

type ChannelsQueue struct {
	channels []*repository.Channel
	lock     sync.Mutex
}

func New() *ChannelsQueue {
	return &ChannelsQueue{}
}

func (q *ChannelsQueue) Enqueue(channel *repository.Channel) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.channels = append(q.channels, channel)
}

func (q *ChannelsQueue) Dequeue() *repository.Channel {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.channels) == 0 {
		return nil
	}

	channel := q.channels[0]
	q.channels = q.channels[1:]

	return channel
}
