package manga

import (
	"github.com/gin-gonic/gin"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetMaidMyHome(c *gin.Context) {
	result := sScrapper.ScrapMaidMyHomePage()

	http_req.RenderResponse(c, 200, result)
}

func GetMaidMySearch(c *gin.Context) {
	query := c.Request.URL.Query().Get("query")

	result := sScrapper.ScrapMaidMyMangaSearchPage(query)

	http_req.RenderResponse(c, 200, result)
}

func GetMaidMyMangaDetail(c *gin.Context) {
	manga_title := c.Request.URL.Query().Get("manga_title")

	result := sScrapper.ScrapMaidMyMangaDetailPage(manga_title)

	http_req.RenderResponse(c, 200, result)
}

func GetMaidMyMangaChapterDetail(c *gin.Context) {
	manga_title := c.Request.URL.Query().Get("manga_chapter")
	manga_chapter := c.Request.URL.Query().Get("manga_chapter")

	result := sScrapper.ScrapMaidMyMangaChapterDetailPage(manga_title, manga_chapter)

	http_req.RenderResponse(c, 200, result)
}
