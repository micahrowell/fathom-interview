package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	// "github.com/gorilla/websocket"
	"github.com/micahrowell/fathom-interview/server/pubsub"
)

var ps = pubsub.NewPubSub()

// var upgrader = websocket.Upgrader{}

func initializer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Connected to server\n")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.Header.Get("topic")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error encountered: %e", err)
		w.Write([]byte("Invalid body"))
		return
	}
	bodyString := string(bodyBytes)
	fmt.Printf("Publishing to topic: %s\nmessage: %s\n", topic, bodyString)
	ps.Publish(topic, bodyString)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.Header.Get("topic")
	ch := ps.Subscribe(topic)
	fmt.Printf("Subscribing to topic: %s\n", topic)
	for msg := range ch {
		fmt.Printf("Sending message: %s\n", msg)
		w.Write([]byte(msg))
	}
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.Header.Get("topic")
	unsubbed := ps.Unsubscribe(topic)
	fmt.Printf("Topic: %s\nUnsubscribed: %t\n", topic, unsubbed)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", initializer)
	mux.HandleFunc("/publish", publishHandler)
	mux.HandleFunc("/subscribe", subscribeHandler)
	mux.HandleFunc("/unsubscribe", unsubscribeHandler)

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
