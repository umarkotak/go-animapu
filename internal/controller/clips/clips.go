package clips

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	rClips "github.com/umarkotak/go-animapu/internal/repository/clips"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetClips(c *gin.Context) {
	clips := rClips.GetClipsFromFirebase()
	http_req.RenderResponse(c, 200, clips)
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
	http_req.RenderResponse(c, 200, clip)
}
