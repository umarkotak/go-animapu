package manga

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/umarkotak/go-animapu/internal/models"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/pkg/app_cache"
	sOnesignal "github.com/umarkotak/go-animapu/internal/service/onesignal"
)

var mangaHubCDN = "https://img.mghubcdn.com/file/imghub"

// UpdateMangaChapters fetch ;atest manga chapter from mangahub
func UpdateMangaChapters(mangaDB models.MangaDB) models.MangaDB {
	appCache := pkgAppCache.GetAppCache()

	res, found := appCache.Get("update_manga_chapter")
	if found {
		fmt.Println("FETCH FROM APP CACHE")
		return res.(models.MangaDB)
	}

	var updatedMangaTitles []string
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
			updatedMangaTitles = append(updatedMangaTitles, mangaTitle)
		} else {
			mangaData.NewAdded++
			fmt.Println("[NO-UPDATE]", mangaTitle)
		}
		mangaDB.MangaDatas[v] = mangaData
	}

	if len(updatedMangaTitles) > 0 {
		joinedString := strings.Join(updatedMangaTitles[:], ", ")
		joinedString = strings.Replace(joinedString, "-", " ", -1)
		fmt.Println("Updated titles: " + joinedString)

		sOnesignal.SendWebNotification("New chapter update!", joinedString)
	}

	appCache.Set("update_manga_chapter", mangaDB, 5*time.Minute)

	return mangaDB
}

// UpdateMangaChaptersV2 fetch latest manga chapter from mangahub
func UpdateMangaChaptersV2(mangaDB models.MangaDB) models.MangaDB {
	// appCache := pkgAppCache.GetAppCache()

	// res, found := appCache.Get("update_manga_chapter")
	// if found {
	// 	fmt.Println("FETCH FROM APP CACHE")
	// 	return res.(models.MangaDB)
	// }

	var updatedMangaTitles []string
	var keys []string
	for k := range mangaDB.MangaDatas {
		keys = append(keys, k)
	}
	// to do update order based on title ascending
	// sort.Strings(keys)

	var wg sync.WaitGroup

	for _, mangaTitle := range keys {
		wg.Add(1)
		go checkMangaLatestChapter(&wg, mangaTitle, &mangaDB, &updatedMangaTitles)
	}

	wg.Wait()

	if len(updatedMangaTitles) > 0 {
		joinedString := strings.Join(updatedMangaTitles[:], ", ")
		joinedString = strings.Replace(joinedString, "-", " ", -1)
		log.Println("Updated titles: " + joinedString)

		sOnesignal.SendWebNotification("New chapter update!", joinedString)
	}

	// appCache.Set("update_manga_chapter", mangaDB, 5*time.Minute)

	return mangaDB
}

func checkMangaLatestChapter(wg *sync.WaitGroup, mangaTitle string, mangaDB *models.MangaDB, updatedMangaTitles *[]string) {
	defer wg.Done()

	mangaData := mangaDB.MangaDatas[mangaTitle]
	mangaLatestChapter := mangaData.MangaLastChapter
	mangaUpdatedChapter := mangaLatestChapter + 1

	if mangaData.Status == "finished" {
		// fmt.Println(mangaTitle, " Finished!")
		return
	}

	targetPathJPG := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/1.jpg"
	targetPathJPG2 := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/5.jpg"
	targetPathPNG := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/1.png"
	targetPathJPEG := mangaHubCDN + "/" + mangaTitle + "/" + strconv.Itoa(mangaUpdatedChapter) + "/5.jpeg"

	resultJpg, _ := http.Get(targetPathJPG)
	resultJpg2, _ := http.Get(targetPathJPG2)
	resultPng, _ := http.Get(targetPathPNG)
	resultJpeg, _ := http.Get(targetPathJPEG)

	if resultJpg.Status == "200 OK" || resultPng.Status == "200 OK" || resultJpeg.Status == "200 OK" || resultJpg2.Status == "200 OK" {
		mangaData.MangaLastChapter = mangaUpdatedChapter
		mangaData.NewAdded = 1
		fmt.Println("[UPDATED]", mangaTitle, " From: ", mangaLatestChapter, " To: ", mangaUpdatedChapter)
		*updatedMangaTitles = append(*updatedMangaTitles, mangaTitle)
	} else {
		// fmt.Println("[NO-UPDATE]", mangaTitle)
		mangaData.NewAdded++
	}
	mangaDB.MangaDatas[mangaTitle] = mangaData
}
