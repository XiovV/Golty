package queue

import (
	"golty/repository"
	"sync"
)

type ChannelsQueue struct {
	channels []*repository.Channel
	lock     *sync.Mutex
	cond     *sync.Cond
}

func New() *ChannelsQueue {
	mutex := &sync.Mutex{}

	return &ChannelsQueue{
		lock: mutex,
		cond: sync.NewCond(mutex),
	}
}

func (q *ChannelsQueue) Enqueue(channel *repository.Channel) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.channels = append(q.channels, channel)
	q.cond.Signal()
}

func (q *ChannelsQueue) Dequeue() *repository.Channel {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.channels) == 0 {
		q.cond.Wait()
	}

	channel := q.channels[0]
	q.channels = q.channels[1:]

	return channel
}
