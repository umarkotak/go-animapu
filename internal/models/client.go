package models

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	firebaseHelper "github.com/umarkotak/go-animapu/internal/pkg/firebase_helper"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var messageObj Message
		json.Unmarshal(message, &messageObj)
		timestampUnix := int32(time.Now().Unix())
		timestampUTC := time.Now().Format("2006-01-02T15:04:05Z07:00")
		messageObj.TimestampUnix = timestampUnix
		messageObj.TimestampUTC = timestampUTC
		processedMessage, err := json.Marshal(messageObj)
		go SetChatMessageToFirebase(messageObj)

		// c.Hub.Broadcast <- message
		c.Hub.Broadcast <- processedMessage
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var ctx = context.Background()

// SetChatMessageToFirebase save message to firebase
func SetChatMessageToFirebase(message Message) Message {
	firebaseDB := firebaseHelper.GetFirebaseDB()

	ref := firebaseDB.NewRef("")
	messagesRef := ref.Child("chat_message_db")
	messageDataRef := messagesRef.Child(strconv.Itoa(int(message.TimestampUnix)))
	messageDataRef.Set(ctx, &message)

	return message
}

func GetChatMessagesFromFirebase() map[string]Message {
	firebaseDB := firebaseHelper.GetFirebaseDB()
	var messages map[string]Message

	ref := firebaseDB.NewRef("")
	messagesRef := ref.Child("chat_message_db")
	if err := messagesRef.Get(ctx, &messages); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	return messages
}
