package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerInfo struct {
	PosX int64
	PosY int64
}

type Player struct {
	Name       string          `json:"name"`
	World      *World          `json:"-"`
	Conn       *websocket.Conn `json:"-"`
	Send       chan []byte     `json:"-"`
	PlayerInfo PlayerInfo      `json:"-"`
}

// ReadPump used to receive message from web socket
func (p *Player) ReadPump() {
	defer func() {
		p.World.Unregister <- p
		p.Conn.Close()
	}()

	p.Conn.SetReadLimit(maxMessageSize)
	p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	p.Conn.SetPongHandler(func(string) error { p.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, rawMessage, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var gameMessage GameMessage
		json.Unmarshal(rawMessage, &gameMessage)
		gameMessage.Player = *p
		gameMessageByte, _ := json.Marshal(gameMessage)

		p.World.Broadcast <- gameMessageByte
	}
}

func (p *Player) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.Send:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(p.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-p.Send)
			}

			err = w.Close()
			if err != nil {
				return
			}

		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
