package pubsub

import (
	"sync"
)

type PubSub interface {
	Publish(topic string, data interface{})
	Subscribe(topic string) <-chan string
	Close()
}

type pubsub struct {
	mu            sync.RWMutex
	subscriptions map[string][]chan string
	closed        bool
}

func NewPubSub() *pubsub {
	ps := &pubsub{}
	ps.subscriptions = map[string][]chan string{}
	return ps
}

func (ps *pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subscriptions[topic] = append(ps.subscriptions[topic], ch)
	return ch
}

func (ps *pubsub) Publish(topic string, data string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subscriptions[topic] {
		go func(channel chan string) {
			channel <- data
		}(ch)
	}
}
