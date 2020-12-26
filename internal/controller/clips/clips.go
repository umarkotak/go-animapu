package clips

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	rClips "github.com/umarkotak/go-animapu/internal/repository/clips"
)

func GetClips(c *gin.Context) {
	clips := rClips.GetClipsFromFirebase()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, clips)
}

func CreateClip(c *gin.Context) {
	type RequestParams struct {
		Content string `json:"content"`
	}
	var json RequestParams
	c.BindJSON(&json)

	clip := models.Clip{
		ID:            strconv.Itoa(int(time.Now().Unix())),
		Content:       json.Content,
		TimestampUnix: int32(time.Now().Unix()),
		TimestampUtc:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	rClips.CreateClipsToFirebase(clip)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, clip)
}
