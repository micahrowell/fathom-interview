package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/micahrowell/fathom-interview/internal"
)

type message internal.Message

func getInput(input chan message, user string) {
	msg := message{}
	in := bufio.NewReader(os.Stdin)
	terminalInput, err := in.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	msg.UserID = user
	msg.Body = terminalInput
	input <- msg
}

func main() {
	var config internal.Configuration
	config.ReadConfig()

	server := fmt.Sprintf("%s:%d", config.Server.Name, config.Server.Port)

	if len(os.Args) != 3 {
		fmt.Println("Client needs Topic and Username arguments")
		return
	}
	topic := os.Args[1]
	username := os.Args[2]
	fmt.Println("Connecting to:", server, "at", topic)

	// create the websocket that will connect to the server
	URL := url.URL{Scheme: "ws", Host: server, Path: topic}
	ws, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Printf("Error creating the websocket: %s", err)
		return
	}
	defer ws.Close()

	// create a channel to listen for this client closing
	// in this example, typing ctrl+C on the keyboard will send the appropriate signal
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt)

	// create a channel to listen for user input from the terminal
	userInputChannel := make(chan message, 1)
	go getInput(userInputChannel, username)

	// create another channel for when we receive an error from the server
	serverInterruptChannel := make(chan struct{})
	go func() {
		defer close(serverInterruptChannel)
		for {
			_, remoteMessage, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error when attempting to read a message from the server: %s\n", err)
				return
			}
			msg := message{}
			err = json.Unmarshal(remoteMessage, &msg)
			if err != nil {
				log.Printf("Error when attempting to unmarshal remote message. Err: %s\nMessage: %s\n", err, string(remoteMessage))
			}
			// display remote message in the termal so the user can see it
			if msg.UserID != username {
				log.Printf("%s: %s", msg.UserID, msg.Body)
			}
		}
	}()

	for {
		select {
		// this case checks for the closure of this channel
		case <-serverInterruptChannel:
			log.Println("Websocket connection closed on server side, exiting...")
			return
		case msg := <-userInputChannel:
			msgJson, err := json.Marshal(msg)
			err = ws.WriteMessage(websocket.TextMessage, msgJson)
			if err != nil {
				log.Println("Write error:", err)
				return
			}
			go getInput(userInputChannel, username)
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
