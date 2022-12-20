package pubsub

import (
	// "fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type PubSub interface {
	Publish(string, int, []byte)
	Subscribe(string, *websocket.Conn)
	Unsubscribe(string, *websocket.Conn) bool
}

type PubSubImpl struct {
	mu            sync.RWMutex
	subscriptions map[string][]*websocket.Conn
}

func NewPubSub() *PubSubImpl {
	ps := &PubSubImpl{}
	ps.subscriptions = map[string][]*websocket.Conn{}
	return ps
}

func (ps *PubSubImpl) Subscribe(topic string, conn *websocket.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subscriptions[topic] = append(ps.subscriptions[topic], conn)
}

func (ps *PubSubImpl) Unsubscribe(topic string, conn *websocket.Conn) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	idx := -1
	for i, c := range ps.subscriptions[topic] {
		if conn.RemoteAddr().Network() == c.RemoteAddr().Network() {
			c.Close()
			idx = i
			break
		}
	}

	if idx > -1 {
		// ps.subscriptions[topic] = append(ps.subscriptions[topic][:idx], ps.subscriptions[topic][idx+1:]...)
		ps.subscriptions[topic][idx] = ps.subscriptions[topic][len(ps.subscriptions[topic])-1]
		ps.subscriptions[topic] = ps.subscriptions[topic][:len(ps.subscriptions[topic])-1]
	}

	return idx != -1
}

func (ps *PubSubImpl) Publish(topic string, messageType int, data []byte) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var err error = nil
	for _, conn := range ps.subscriptions[topic] {
		err = conn.WriteMessage(messageType, data)
	}
	return err
}
