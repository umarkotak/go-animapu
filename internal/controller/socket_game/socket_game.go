package socket_game

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/umarkotak/go-animapu/internal/models"
	sUser "github.com/umarkotak/go-animapu/internal/service/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Serve(c *gin.Context, world *models.World) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	auth := c.Query("token")

	userData, err := sUser.DetailService(auth)
	if err != nil {
		renderError(c, err)
		return
	}

	player := &models.Player{
		Name:  userData.Username,
		World: world,
		Conn:  conn,
		Send:  make(chan []byte, 256),
		PlayerInfo: models.PlayerInfo{
			PosX: 1,
			PosY: 1,
		},
	}
	player.World.Register <- player

	go player.WritePump()
	go player.ReadPump()
}

func renderError(c *gin.Context, err error) {
	var response gin.H
	var statusCode int

	response = gin.H{
		"message": err.Error(),
	}
	statusCode = 400

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(statusCode, response)
}
