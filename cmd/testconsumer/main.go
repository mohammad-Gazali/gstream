package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)

	if err != nil {
		log.Fatal(err)
	}

	msg := WSMessage{
		Action: "subscribe",
		Topics:  []string{"foobaz-1", "foobaz-2"},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	// if we don't add this for loop then the connection will end immediately
	for {
		var msg WSMessage

		// here the ReadJSON will block until we receive message from the server
		// this is because we use websockets
		if err := conn.ReadJSON(&msg); err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	}
}
