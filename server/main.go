package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/micahrowell/fathom-interview/server/pubsub"
)

var ps = pubsub.NewPubSub()

func initializer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Connected to server\n")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.Header.Get("topic")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println("Server up and listening on port 3000...")
	if err != nil {
		panic(err)
	}
}
