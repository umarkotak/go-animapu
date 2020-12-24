package chats

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/umarkotak/go-animapu/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// GetChatsConnection start chat websocket
func GetChatsConnection(c *gin.Context, hub *models.Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &models.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
