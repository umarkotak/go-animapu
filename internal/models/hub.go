package models

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			execAfterRegister(client)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func execAfterRegister(client *Client) {
	data := fmt.Sprintf(
		`
			{
				"type": "normal_message",
				"timestamp_unix": 0,
				"name": "%v",
				"message": "welcome to chatto_v1"
			}
		`,
		"animapu",
	)
	client.Send <- []byte(data)

	messages := GetChatMessagesFromFirebase()
	var arrMessages []Message

	for _, message := range messages {
		arrMessages = append(arrMessages, message)
	}

	sort.Slice(arrMessages, func(i, j int) bool {
		return arrMessages[i].TimestampUTC < arrMessages[j].TimestampUTC
	})

	var byteMessage []byte
	for _, sMessage := range arrMessages {
		byteMessage, _ = json.Marshal(sMessage)
		client.Send <- []byte(byteMessage)
		fmt.Println("Sent! ", sMessage)
	}

}
