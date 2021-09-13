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
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	sStatistic "github.com/umarkotak/go-animapu/internal/service/statistic"
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
		mangaDB = sManga.UpdateMangaChaptersV2(mangaDB)
		go rManga.UpdateMangaToFireBase(mangaDB)
		appCache.GetAppCache().Set("UPDATED_MANGA_CACHE", "UPDATED", 1*time.Minute)
	}

	http_req.RenderResponse(c, 200, mangaDB)
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

func GetMangaStatistics(c *gin.Context) {
	result := sStatistic.GenerateMangaStatistic()

	http_req.RenderResponse(c, 200, result)
}

func GetDailyMangaStatistics(c *gin.Context) {
	result := sStatistic.GenerateDailyMangaStatistic()

	http_req.RenderResponse(c, 200, result)
}

func GetMangaDetail(c *gin.Context) {
	manga_title := c.Request.URL.Query().Get("manga_title")

	result := sScrapper.GetMangaDetailV1(manga_title)

	http_req.RenderResponse(c, 200, result)
}

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
