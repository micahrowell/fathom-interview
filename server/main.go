package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/micahrowell/fathom-interview/server/pubsub"
)

// TODO: modify to remove connection, publish to all connections
func listen(conn *websocket.Conn, ps pubsub.PubSub, path string) {
	ps.Subscribe(path)
	for {
		// read a message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message
		fmt.Println(string(messageContent))

		// reponse message
		messageResponse := fmt.Sprintf("Your message is: %s\n", messageContent)
		ps.Publish(path, messageResponse)

		if err := conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			log.Println(err)
			return
		}
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
		listen(ws, ps, r.URL.Path)
	})

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
