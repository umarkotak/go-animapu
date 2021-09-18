package manga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/umarkotak/go-animapu/internal/models"
	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
	appCache "github.com/umarkotak/go-animapu/internal/utils/app_cache"
	"github.com/umarkotak/go-animapu/internal/utils/http_req"
)

// GetManga get list of all manga in DB
func GetManga(c *gin.Context) {
	jsonFile, err := os.Open("data/mangas.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	http_req.RenderResponse(c, 200, result)
}

// UpdateManga run update manga
func UpdateManga(c *gin.Context) {
	mangaDB := rManga.GetMangaFromJSON()
	mangaDB = sManga.UpdateMangaChapters(mangaDB)
	mangaDB = rManga.UpdateMangaToJSON(mangaDB)

	http_req.RenderResponse(c, 200, mangaDB)
}

// GetMangaFirebase get manga from firebase
func GetMangaFirebase(c *gin.Context) {
	// mangaDB := rManga.GetMangaFromFireBaseV2()
	mangaDB := rManga.GetMangaFromFireBaseV2WithoutCache()

	http_req.RenderResponse(c, 200, mangaDB)
}

// UpdateMangaFirebase update mangat to firebase
func UpdateMangaFirebase(c *gin.Context) {
	updated, _ := appCache.GetAppCache().Get("UPDATED_MANGA_CACHE")

	// mangaDB := rManga.GetMangaFromFireBaseV2()
	mangaDB := rManga.GetMangaFromFireBaseV2WithoutCache()
	// mangaDB = sManga.UpdateMangaChapters(mangaDB)
	if updated == nil {
		fmt.Println("DIRECT UPDATE")
		appCache.GetAppCache().Set("UPDATED_MANGA_CACHE", "UPDATED", 3*time.Minute)
		mangaDB = sManga.UpdateMangaChaptersV2(mangaDB)
		go rManga.UpdateMangaToFireBase(mangaDB)
	}

	http_req.RenderResponse(c, 200, mangaDB)
}

func PostAddToGeneralMangaLibrary(c *gin.Context) {
	mangaData := models.MangaData{
		Title:            "",
		CompactTitle:     "",
		MangaLastChapter: 0,
		AveragePage:      100,
		Status:           "ongoing",
		ImageURL:         "",
		NewAdded:         1,
		Weight:           10000,
		Finder:           "EXTERNAL",
	}
	c.BindJSON(&mangaData)
	mangaData, _ = rManga.AddMangaToFireBaseGeneralLibrary(mangaData)

	http_req.RenderResponse(c, 200, mangaData)
}
