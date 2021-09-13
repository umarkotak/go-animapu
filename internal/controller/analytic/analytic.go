package analytic

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	sAnalytic "github.com/umarkotak/go-animapu/internal/service/analytic"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
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

	go sAnalytic.LogUserAnalytic(requestParams.MangaTitle, requestParams.MangaPage, requestParams.UserIP)

	go sAnalytic.LogDailyMangaView(requestParams.MangaTitle, requestParams.MangaPage, requestParams.UserIP)

	http_req.RenderResponse(c, 200, "OK")
}

func sanitizeClientIP(clientIP string) string {
	sanitizedIP := clientIP
	sanitizedIP = strings.Replace(sanitizedIP, "[", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "]", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, ".", "-", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "/", "", -1)
	sanitizedIP = strings.Replace(sanitizedIP, "$", "", -1)

	return sanitizedIP
}
