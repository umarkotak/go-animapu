package manga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"

	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
	sScrapper "github.com/umarkotak/go-animapu/internal/service/scrapper"
	sStatistic "github.com/umarkotak/go-animapu/internal/service/statistic"
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

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, result)
}

// UpdateManga run update manga
func UpdateManga(c *gin.Context) {
	mangaDB := rManga.GetMangaFromJSON()
	mangaDB = sManga.UpdateMangaChapters(mangaDB)
	mangaDB = rManga.UpdateMangaToJSON(mangaDB)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, mangaDB)
}

// GetMangaFirebase get manga from firebase
func GetMangaFirebase(c *gin.Context) {
	mangaDB := rManga.GetMangaFromFireBaseV2()

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, mangaDB)
}

// UpdateMangaFirebase update mangat to firebase
func UpdateMangaFirebase(c *gin.Context) {
	mangaDB := rManga.GetMangaFromFireBaseV2()
	mangaDB = sManga.UpdateMangaChapters(mangaDB)
	mangaDB = rManga.UpdateMangaToFireBase(mangaDB)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, mangaDB)
}

// GetMangaSearch search manga title
func GetMangaSearch(c *gin.Context) {
	title := c.Query("title")

	mangaDB := sScrapper.SearchMangaTitle(title)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, mangaDB)
}

// GetMangaTodays list of todays manga
func GetMangaTodays(c *gin.Context) {
	mangaDB := sScrapper.GetTodaysMangaTitleV2()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, mangaDB)
}

func GetMangaStatistics(c *gin.Context) {
	result := sStatistic.GenerateMangaStatistic()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(200, result)
}
