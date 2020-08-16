package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/umarkotak/go-animapu/internal/models"
)

// GetMangaFromJSON read manga from data/
func GetMangaFromJSON() {
	var mangaDbFilePath = "../../data/mangas_play.json"
	mangaDbJSONFile, err := os.Open(mangaDbFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer mangaDbJSONFile.Close()

	mangaDbByteValue, _ := ioutil.ReadAll(mangaDbJSONFile)

	var mangaDB models.MangaDB
	json.Unmarshal([]byte(mangaDbByteValue), &mangaDB)

	return mangaDB
}
