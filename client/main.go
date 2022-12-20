package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func getInput(input chan string, user string) {
	in := bufio.NewReader(os.Stdin)
	terminalInput, err := in.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	message := fmt.Sprintf("%s: %s", user, terminalInput)
	input <- message
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Client needs Server, Path, and Username arguments")
		return
	}
	server := os.Args[1]
	path := os.Args[2]
	username := os.Args[3]
	fmt.Println("Connecting to:", server, "at", path)

	// create a channel to listen for this client closing
	// in this example, typing ctrl+C on the keyboard will send the appropriate signal
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt)

	// create a channel to listen for user input in terminal
	inputChannel := make(chan string, 1)
	go getInput(inputChannel, username)

	// create the websocket that will connect to the server
	URL := url.URL{Scheme: "ws", Host: server, Path: path}
	ws, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Printf("Error creating the websocket: %e", err)
		return
	}
	defer ws.Close()

	// create another channel for when we receive an error from the server
	serverInterruptChannel := make(chan struct{})
	go func() {
		defer close(serverInterruptChannel)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error when attempting to read a message from the server: %s\n", err)
				return
			}
			log.Println(string(message))
		}
	}()

	for {
		select {
		// this case checks for the closure of this channel
		case <-serverInterruptChannel:
			log.Println("Websocket connection closed on server side, exiting...")
			return
		case msg := <-inputChannel:
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Write error:", err)
				return
			}
			go getInput(inputChannel, username)
		case <-quitChannel:
			log.Println("Client is quitting!")
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Printf("Error during closing websocket: %s", err)
				return
			}
			return
		}
	}
}
