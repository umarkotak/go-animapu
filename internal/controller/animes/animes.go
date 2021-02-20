package animes

import (
	"github.com/gin-gonic/gin"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
)

func GetAnimesMap(c *gin.Context) {
	result := sScrapper.FetchAllAnime(c)
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, result)
}
