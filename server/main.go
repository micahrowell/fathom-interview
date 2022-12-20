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

func subscribeAndListen(conn *websocket.Conn, ps *pubsub.PubSub, topic string) {
	ps.Subscribe(topic, conn)

	// listen for messages
	for {
		// read incoming message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// convert from JSON
		msg := internal.Message{}
		err = json.Unmarshal(messageContent, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		// display message on the server console
		fmt.Printf("%s - %s: %s", topic, msg.UserID, msg.Body)

		// send the message to all subscribers
		_ = ps.Publish(topic, messageType, messageContent)
	}
}

func main() {
	var config internal.Configuration
	config.ReadConfig()

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
