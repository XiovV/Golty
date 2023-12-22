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

	for _, c := range q.channels {
		if c.ID == channel.ID {
			return
		}
	}

	q.channels = append(q.channels, channel)
	q.cond.Signal()
}

func (q *ChannelsQueue) Dequeue() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.channels) == 0 {
		q.cond.Wait()
	}

	q.channels = q.channels[1:]
}

func (q *ChannelsQueue) GetFirst() *repository.Channel {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.channels) == 0 {
		q.cond.Wait()
	}

	return q.channels[0]
}
