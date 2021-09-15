package manga

import (
	"github.com/gin-gonic/gin"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetMangaDetail(c *gin.Context) {
	manga_title := c.Request.URL.Query().Get("manga_title")

	result := sScrapper.GetMangaDetailV1(manga_title)

	http_req.RenderResponse(c, 200, result)
}

// GetMangaSearch search manga title
func GetMangaSearch(c *gin.Context) {
	title := c.Query("title")

	mangaDB := sScrapper.SearchMangaTitle(title)

	http_req.RenderResponse(c, 200, mangaDB)
}

// GetMangaTodays list of todays manga
func GetMangaTodays(c *gin.Context) {
	mangaDB := sScrapper.GetTodaysMangaTitleV2()

	http_req.RenderResponse(c, 200, mangaDB)
}
