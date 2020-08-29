package manga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/umarkotak/go-animapu/internal/models"
)

var mangaDbFilePath = "data/mangas.json"

// var mangaDbFilePath = "data/mangas.json"

// GetMangaFromJSON read manga from json db in data/
func GetMangaFromJSON() models.MangaDB {
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

// UpdateMangaToJSON update manga to json db in data/
func UpdateMangaToJSON(mangaDb models.MangaDB) models.MangaDB {
	mangaDbJSON, _ := json.Marshal(mangaDb)
	ioutil.WriteFile(mangaDbFilePath, mangaDbJSON, 0644)

	return mangaDb
}

// GetMangaFromFireBase fetch manga from firebase
func GetMangaFromFireBase() {

}
