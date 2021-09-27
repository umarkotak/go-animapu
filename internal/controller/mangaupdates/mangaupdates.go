package mangaupdates

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetSeries(c *gin.Context) {
	mangaID := c.Query("manga_id")
	result, err := scrapper.MangaupdatesGetSeries(mangaID)
	if err != nil {
		http_req.RenderResponse(c, 422, err)
	}

	http_req.RenderResponse(c, 200, result)
}

func GetReleases(c *gin.Context) {
	mangaDB, err := scrapper.MangaupdatesGetReleases()
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}

	http_req.RenderResponse(c, 200, mangaDB)
}

func GetDetailByTitle(c *gin.Context) {
	mangaTitle := c.Param("manga_title")
	mangaDetail, err := scrapper.MangaupdatesSeriesDetailByTitle(mangaTitle)
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}

	http_req.RenderResponse(c, 200, mangaDetail)
}

func Search(c *gin.Context) {
	title := c.Query("title")

	mangaDB, err := scrapper.MangaupdatesSearch(title)
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}

	http_req.RenderResponse(c, 200, mangaDB)
}

func ReleasesSearch(c *gin.Context) {
	mangaupdateID := c.Param("mangaupdates_id")

	mangaDetail, err := scrapper.MangaupdatesSeriesDetailByID(mangaupdateID)
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}

	http_req.RenderResponse(c, 200, mangaDetail)
}
