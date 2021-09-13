package animes

import (
	"github.com/gin-gonic/gin"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetAnimesMap(c *gin.Context) {
	result := sScrapper.FetchAllAnime(c)
	http_req.RenderResponse(c, 200, result)
}

func GetSearchAnime(c *gin.Context) {
	result := sScrapper.SearchAnime(c)
	http_req.RenderResponse(c, 200, result)
}
