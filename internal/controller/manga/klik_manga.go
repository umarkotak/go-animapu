package manga

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	"github.com/umarkotak/go-animapu/internal/repository/manga"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetKlikMangaHome(c *gin.Context) {
	result, err := manga.GetMangaDBFromCache("GetKlikMangaHome")
	if err == nil {
		http_req.RenderResponse(c, 200, result)
		return
	}

	result = sScrapper.ScrapKlikMangaHomePage()
	if len(result.MangaDataKeys) > 0 {
		manga.SetMangaDBToCache("GetKlikMangaHome", result)
	}

	http_req.RenderResponse(c, 200, result)
}

func GetKlikMangaHomeNextPage(c *gin.Context) {
	page := c.Param("page")
	pageInt64, _ := strconv.ParseInt(page, 10, 64)

	result := sScrapper.ScrapKlikMangaHomeNextPage(pageInt64)

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

func GetKlikMangaSearchNextPage(c *gin.Context) {
	searchParams := models.KlikMangaSearchParams{
		Title:  c.Request.URL.Query().Get("title"),
		Status: c.Request.URL.Query().Get("status"),
		Genre:  c.Request.URL.Query().Get("genre"),
		Page:   c.Param("page"),
	}

	result := sScrapper.ScrapKlikMangaSearchNextPage(searchParams)

	http_req.RenderResponse(c, 200, result)
}
