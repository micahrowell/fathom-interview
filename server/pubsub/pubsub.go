package pubsub

import (
	"sync"
)

type PubSub interface {
	Publish(topic string, data string)
	Subscribe(topic string) <-chan string
	Unsubscribe(topic string) bool
	// Close()
}

type PubSubImpl struct {
	mu            sync.RWMutex
	subscriptions map[string][]chan string
	closed        bool
}

func NewPubSub() *PubSubImpl {
	ps := &PubSubImpl{}
	ps.subscriptions = map[string][]chan string{}
	return ps
}

func (ps *PubSubImpl) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subscriptions[topic] = append(ps.subscriptions[topic], ch)
	return ch
}

func (ps *PubSubImpl) Unsubscribe(topic string) bool {
	unsubbed := false
	for _, c := range ps.subscriptions[topic] {
		close(c)
		unsubbed = true
	}

	if unsubbed {
		delete(ps.subscriptions, topic)
	}

	return unsubbed
}

func (ps *PubSubImpl) Publish(topic string, data string) {
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
