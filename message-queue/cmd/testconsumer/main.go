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

type Message struct {
	// 어떤 토픽으로 쏠지
	Topic string `json:"topic"`
	// 어떤 데이터를 담고있는지
	Data []byte `json:"data"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer conn.Close()

	msg := WSMessage{
		Action: "subscribe",
		Topics: []string{"foobarbaz"},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Fatal("WriteJSON:", err)
	}

	for {
		msg := Message{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ReadJSON:", err)
			break
		}
		fmt.Printf("Received message: Topic: %s, Data: %s\n", msg.Topic, string(msg.Data))
	}
}
