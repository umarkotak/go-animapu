package analytic

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	sAnalytic "github.com/umarkotak/go-animapu/internal/service/analytic"
)

func PostUserAnalyticV1(c *gin.Context) {
	type RequestParams struct {
		MangaTitle string `json:"MangaTitle"`
		MangaPage  int    `json:"MangaPage"`
		UserIP     string
	}
	var requestParams RequestParams
	c.BindJSON(&requestParams)
	requestParams.UserIP = sanitizeClientIP(c.ClientIP())
	fmt.Println("DEBUG MSG", requestParams.UserIP)

	sAnalytic.LogUserAnalytic(requestParams.MangaTitle, requestParams.MangaPage, requestParams.UserIP)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, "OK")
}

// SkipCors skip cors
func SkipCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, "CORS OK")
}

func sanitizeClientIP(clientIP string) string {
	sanitizedIP := clientIP
	sanitizedIP = strings.Replace(sanitizedIP, "[", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "]", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, ".", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "/", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "$", "", -1)

	return sanitizedIP
}
