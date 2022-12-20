package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type PubSub interface {
	Publish(string, int, string)
	Subscribe(string, *websocket.Conn)
	// Unsubscribe(string, *websocket.Conn) bool
	// Close()
}

type PubSubImpl struct {
	mu            sync.RWMutex
	subscriptions map[string][]*websocket.Conn
	closed        bool
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

// func (ps *PubSubImpl) Unsubscribe(topic string, conn *websocket.Conn) bool {
// 	unsubbed := false
// 	for _, c := range ps.subscriptions[topic] {
// 		c.Close()
// 		unsubbed = true
// 	}

// 	if unsubbed {
// 		// delete(ps.subscriptions, topic)
// 	}

// 	return unsubbed
// }

func (ps *PubSubImpl) Publish(topic string, messageType int, data string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil
	}

	var err error = nil
	for _, conn := range ps.subscriptions[topic] {
		err = conn.WriteMessage(messageType, []byte(data))
	}
	return err
}
