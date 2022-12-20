package main

import (
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
	ps.Publish(topic, bodyString)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.Header.Get("topic")
	ch := ps.Subscribe(topic)
	for msg := range ch {
		w.Write([]byte(msg))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", initializer)
	mux.HandleFunc("/publish", publishHandler)
	mux.HandleFunc("/subscribe", subscribeHandler)

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		panic(err)
	}
}
