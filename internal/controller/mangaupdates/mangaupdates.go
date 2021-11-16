package mangaupdates

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/go-animapu/internal/models"
	"github.com/umarkotak/go-animapu/internal/repository/manga"
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
	mangaDB, err := manga.GetMangaDBFromCache("MangaupdatesGetReleases")
	if err == nil {
		http_req.RenderResponse(c, 200, mangaDB)
		return
	}

	mangaDB, err = scrapper.MangaupdatesGetReleases()
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}
	if len(mangaDB.MangaDatas) > 0 {
		manga.SetMangaDBToCache("MangaupdatesGetReleases", mangaDB)
	}

	http_req.RenderResponse(c, 200, mangaDB)
}

func GetReleasesV2(c *gin.Context) {
	animapuMangas := []models.MangaData{}

	err := manga.GetAnimapuMangasFromCache("MangaupdatesGetReleasesV2", &animapuMangas)
	if err == nil {
		http_req.RenderResponse(c, 200, animapuMangas)
		return
	}

	animapuMangas, err = scrapper.MangaupdatesGetReleasesV2()
	if err != nil {
		http_req.RenderResponse(c, 422, err)
		return
	}

	if len(animapuMangas) > 0 {
		manga.SetObjectToCache("MangaupdatesGetReleasesV2", animapuMangas)
	}

	http_req.RenderResponse(c, 200, animapuMangas)
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
