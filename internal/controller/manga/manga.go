package manga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"

	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
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
