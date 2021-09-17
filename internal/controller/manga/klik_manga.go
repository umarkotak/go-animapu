package manga

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetKlikMangaHome(c *gin.Context) {
	result := sScrapper.ScrapKlikMangaHomePage()

	http_req.RenderResponse(c, 200, result)
}

func GetKlikMangaDetail(c *gin.Context) {
	manga_title := c.Param("manga_id")

	result := sScrapper.ScrapKlikMangaDetailPage(manga_title)

	http_req.RenderResponse(c, 200, result)
}

func GetKlikMangaChapterDetail(c *gin.Context) {
	manga_title := c.Param("manga_id")
	manga_chapter := c.Param("manga_chapter")

	result := sScrapper.ScrapKlikMangaChapterDetailPage(manga_title, manga_chapter)

	http_req.RenderResponse(c, 200, result)
}

func GetKlikMangaSearch(c *gin.Context) {
	searchParams := models.KlikMangaSearchParams{
		Title:  c.Request.URL.Query().Get("title"),
		Status: c.Request.URL.Query().Get("status"),
		Genre:  c.Request.URL.Query().Get("genre"),
	}

	result := sScrapper.ScrapKlikMangaSearch(searchParams)

	http_req.RenderResponse(c, 200, result)
}
