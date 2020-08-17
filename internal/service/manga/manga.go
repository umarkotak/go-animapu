package manga

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/umarkotak/go-animapu/internal/models"
)

var mangaHubCDN = "https://img.mghubcdn.com/file/imghub"
var respJpg, respPng *http.Response

// UpdateMangaChapters fetch ;atest manga chapter from mangahub
func UpdateMangaChapters(mangaDB models.MangaDB) models.MangaDB {
	var keys []string
	for k := range mangaDB.MangaDatas {
		keys = append(keys, k)
	}
	// to do update order based on title ascending
	sort.Strings(keys)
	for _, v := range keys {
		mangaTitle := v
		mangaData := mangaDB.MangaDatas[v]
		mangaLatestChapter := mangaData.MangaLastChapter
		mangaUpdatedChapter := mangaLatestChapter + 1

		if mangaData.Status == "finished" {
			fmt.Println(mangaTitle, " Finished!")
			continue
		}

		targetPathJPG := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/1.jpg"
		targetPathPNG := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/1.png"

		resultJpg, _ := http.Get(targetPathJPG)
		resultPng, _ := http.Get(targetPathPNG)

		if resultJpg.Status == "200 OK" || resultPng.Status == "200 OK" {
			mangaData.MangaLastChapter = mangaUpdatedChapter
			mangaData.NewAdded = 1
			fmt.Println("[UPDATED]", mangaTitle, " From: ", mangaLatestChapter, " To: ", mangaUpdatedChapter)
		} else {
			mangaData.NewAdded++
			fmt.Println("[NO-UPDATE]", mangaTitle)
		}
		mangaDB.MangaDatas[v] = mangaData
	}

	return mangaDB
}
