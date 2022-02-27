package manga

import (
	"strconv"

	"github.com/gin-gonic/gin"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

func GetMangabatHome(c *gin.Context) {
	// result, err := manga.GetMangaDBFromCache("GetKlikMangaHome")
	// if err == nil {
	// 	http_req.RenderResponse(c, 200, result)
	// 	return
	// }

	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	if page <= 0 {
		page = 1
	}

	result := sScrapper.ScrapMangabatList(page)
	// if len(result.MangaDataKeys) > 0 {
	// 	manga.SetMangaDBToCache("GetKlikMangaHome", result)
	// }

	http_req.RenderResponse(c, 200, result)
}

func GetMangabatMangaDetail(c *gin.Context) {
	mangaID := c.Param("manga_id")

	result := sScrapper.ScrapMangabatDetail(mangaID)

	http_req.RenderResponse(c, 200, result)
}

func GetMangabatMangaChapterDetail(c *gin.Context) {
	mangaID := c.Param("manga_id")
	chapterID := c.Param("chapter_id")

	result := sScrapper.ScrapMangabatChapterDetail(mangaID, chapterID)

	http_req.RenderResponse(c, 200, result)
}
