package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/micahrowell/fathom-interview/internal"
	"github.com/micahrowell/fathom-interview/server/pubsub"

	"github.com/gorilla/websocket"
)

// TODO: modify to remove connection, publish to all connections
func subscribeAndListen(conn *websocket.Conn, ps *pubsub.PubSubImpl, path string) {
	ps.Subscribe(path, conn)

	// listen for messages
	for {
		// read a message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		msg := internal.Message{}
		err = json.Unmarshal(messageContent, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		// display message on the server console
		messageWithTopic := fmt.Sprintf("%s - %s: %s", path, msg.UserID, msg.Body)
		fmt.Println(messageWithTopic)

		// send the message to all subscribers
		_ = ps.Publish(path, messageType, messageContent)
	}
}

func main() {
	ps := pubsub.NewPubSub()
	var upgrader = websocket.Upgrader{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Websocket connected at path %s", r.URL.Path)
		subscribeAndListen(ws, ps, r.URL.Path)
	})

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
