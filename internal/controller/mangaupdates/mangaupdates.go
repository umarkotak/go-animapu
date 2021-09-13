package mangaupdates

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetSeries(c *gin.Context) {
	mangaID := c.Query("manga_id")
	result, err := scrapper.GetSeries(mangaID)
	if err != nil {
		http_req.RenderResponse(c, 422, err)
	}

	http_req.RenderResponse(c, 200, result)
}

func GetReleases(c *gin.Context) {
	mangaDB, err := scrapper.GetReleases()
	if err != nil {
		http_req.RenderResponse(c, 422, err)
	}

	http_req.RenderResponse(c, 200, mangaDB)
}
