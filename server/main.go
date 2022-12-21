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

const (
	READMSGERR   = "Could not read the incoming message"
	UNMARSHALERR = "Could not unmarshal the JSON"
	RESPONDERR   = "Could not send a response to the client"
	PUBLISHERR   = "Could not publish the message to all subscribers"
)

type message internal.Message

func subscribeAndListen(conn *websocket.Conn, ps *pubsub.PubSub, topic string) {
	ps.Subscribe(topic, conn)

	// listen for messages
	for {
		// read incoming message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			err = sendErrorMessage(conn, messageType, READMSGERR)
			if err != nil {
				log.Println(RESPONDERR)
				log.Println(err)
			}
			continue
		}

		// convert from JSON
		msg := message{}
		err = json.Unmarshal(messageContent, &msg)
		if err != nil {
			log.Println(err)
			err = sendErrorMessage(conn, messageType, UNMARSHALERR)
			if err != nil {
				log.Println(RESPONDERR)
				log.Println(err)
			}
			continue
		}

		// display message on the server console
		fmt.Printf("%s - %s: %s", topic, msg.UserID, msg.Body)

		// send the message to all subscribers
		err = ps.Publish(topic, messageType, messageContent)
		if err != nil {
			log.Println(err)
			err = sendErrorMessage(conn, messageType, PUBLISHERR)
			if err != nil {
				log.Println(RESPONDERR)
				log.Println(err)
			}
		}
	}
}

func sendErrorMessage(conn *websocket.Conn, msgType int, info string) error {
	msg := fmt.Sprintf("The server encountered an error with your message: %s", info)
	return conn.WriteMessage(msgType, []byte(msg))
}

func main() {
	var config internal.Configuration
	config.ReadConfig()

	ps := pubsub.NewPubSub()
	upgrader := websocket.Upgrader{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Websocket connected at %s", r.URL.Path)
		go subscribeAndListen(ws, ps, r.URL.Path)
	})

	port := fmt.Sprintf(":%d", config.Server.Port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
