package manga

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sOnesignal "github.com/umarkotak/go-animapu/internal/service/onesignal"
	"github.com/umarkotak/go-animapu/internal/service/scrapper"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/utils/app_cache"
)

var (
	mangaHubCDN  = "https://img.mghubcdn.com/file/imghub"
	mangaDBMutex = sync.RWMutex{}
)

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

	tempMangaDB := mangaDB

	var updatedMangaTitles []string
	var keys []string
	for k := range mangaDB.MangaDatas {
		keys = append(keys, k)
	}
	// to do update order based on title ascending
	// sort.Strings(keys)

	var wg sync.WaitGroup

	currCount := 1
	limit := 10000

	logrus.Infof("executing update")
	for idx, mangaTitle := range keys {
		if currCount >= limit {
			continue
		}

		wg.Add(1)
		// go checkMangaLatestChapter(&wg, mangaTitle, &tempMangaDB, &updatedMangaTitles)
		// go checkMangaLatestChapterV2(&wg, mangaTitle, &tempMangaDB, &updatedMangaTitles)
		// syncMangaUpdatesID(&wg, mangaTitle, &tempMangaDB, &updatedMangaTitles)
		logrus.Infof("UPDATING %v/%v: %v\n", idx, len(keys), mangaTitle)
		checkMangaLatestChapterV3(&wg, mangaTitle, &tempMangaDB, &updatedMangaTitles, idx)

		currCount++
	}

	wg.Wait()

	// if len(updatedMangaTitles) > 0 {
	// 	joinedString := strings.Join(updatedMangaTitles[:], ", ")
	// 	joinedString = strings.Replace(joinedString, "-", " ", -1)
	// 	log.Println("Updated titles: " + joinedString)

	// 	sOnesignal.SendWebNotification("New chapter update!", joinedString)
	// }

	// appCache.Set("update_manga_chapter", mangaDB, 5*time.Minute)

	return tempMangaDB
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", targetPathJPG, nil)
	req.Header.Add("authority", `img.mghubcdn.com`)
	resp, _ := client.Do(req)
	fmt.Println(targetPathJPG, resp.Status)

	resultJpg, _ := http.Get(targetPathJPG)
	resultJpg2, _ := http.Get(targetPathJPG2)
	resultPng, _ := http.Get(targetPathPNG)
	resultJpeg, _ := http.Get(targetPathJPEG)

	// fmt.Println(targetPathJPG)
	if resultJpg.Status == "200 OK" || resultPng.Status == "200 OK" || resultJpeg.Status == "200 OK" || resultJpg2.Status == "200 OK" {
		mangaData.MangaLastChapter = mangaUpdatedChapter
		mangaData.NewAdded = 1
		fmt.Println("[UPDATED]", mangaTitle, " From: ", mangaLatestChapter, " To: ", mangaUpdatedChapter)
		*updatedMangaTitles = append(*updatedMangaTitles, mangaTitle)
	} else {
		// fmt.Println("[NO-UPDATE]", mangaTitle, resultJpg.Status, resultPng.Status, resultJpeg.Status, resultJpg2.Status)
		mangaData.NewAdded++
	}
	fmt.Println()
	mangaDB.MangaDatas[mangaTitle] = mangaData
}

func checkMangaLatestChapterV2(wg *sync.WaitGroup, mangaTitle string, mangaDB *models.MangaDB, updatedMangaTitles *[]string) {
	defer wg.Done()

	mangaDBMutex.RLock()
	mangaData := mangaDB.MangaDatas[mangaTitle]
	mangaDBMutex.RUnlock()

	mangaData.Title = mangaTitle
	mangaData.CompactTitle = strings.Replace(mangaTitle, "-", " ", -1)

	mangaDBMutex.Lock()
	mangaDB.MangaDatas[mangaTitle] = mangaData
	mangaDBMutex.Unlock()

	if mangaData.Status == "finished" {
		return
	}

	mangaDetail := scrapper.GetMangaDetailV1(mangaTitle)

	if len(mangaDetail.Chapters) == 0 {
		return
	}

	lastChapter, err := strconv.ParseFloat(mangaDetail.Chapters[0], 32)
	if err != nil || lastChapter <= 0 {
		return
	}

	lastChapterRounded := int(lastChapter)

	if lastChapterRounded <= mangaData.MangaLastChapter {
		mangaData.NewAdded++
		return
	}

	mangaData.MangaLastChapter = lastChapterRounded
	mangaData.NewAdded = 1

	mangaDBMutex.Lock()
	mangaDB.MangaDatas[mangaTitle] = mangaData
	mangaDBMutex.Unlock()
}

func syncMangaUpdatesID(wg *sync.WaitGroup, mangaTitle string, mangaDB *models.MangaDB, updatedMangaTitles *[]string) {
	defer wg.Done()

	mangaDBMutex.RLock()
	mangaData := mangaDB.MangaDatas[mangaTitle]
	mangaDBMutex.RUnlock()

	mangaData.Title = mangaTitle
	mangaData.CompactTitle = strings.Replace(mangaTitle, "-", " ", -1)

	mangaDBMutex.Lock()
	mangaDB.MangaDatas[mangaTitle] = mangaData
	mangaDBMutex.Unlock()

	if mangaData.Status == "finished" {
		return
	}

	if mangaData.MangaUpdatesID != "" {
		return
	}

	searchMangaDB, err := scrapper.MangaupdatesSearch(mangaData.CompactTitle)
	if err != nil {
		logrus.Errorln("scrapper.MangaupdatesSearch", err)
	}
	if len(searchMangaDB.MangaDataKeys) == 0 {
		return
	}
	matchTitle := searchMangaDB.MangaDataKeys[0]
	matchMangaData := searchMangaDB.MangaDatas[matchTitle]
	mangaData.MangaUpdatesID = matchMangaData.MangaUpdatesID

	logrus.Infof("%v = %v; %+v\n", mangaTitle, matchMangaData.MangaUpdatesID, searchMangaDB.MangaDataKeys)

	mangaDBMutex.Lock()
	mangaDB.MangaDatas[mangaTitle] = mangaData
	mangaDBMutex.Unlock()
}

func checkMangaLatestChapterV3(wg *sync.WaitGroup, mangaTitle string, mangaDB *models.MangaDB, updatedMangaTitles *[]string, idx int) {
	defer wg.Done()
	// rand.Seed(time.Now().UnixNano())
	// n := rand.Intn(10)
	// time.Sleep(time.Duration(idx*1250) * time.Millisecond)

	mangaDBMutex.RLock()
	mangaData := mangaDB.MangaDatas[mangaTitle]
	mangaDBMutex.RUnlock()

	mangaData.Title = mangaTitle
	mangaData.CompactTitle = strings.Replace(mangaTitle, "-", " ", -1)

	if mangaData.Status == "finished" {
		return
	}

	if mangaData.MangaUpdatesID == "" {
		return
	}

	if mangaData.NewAdded <= 3 {
		mangaData.NewAdded++
		return
	}

	if time.Now().Unix()-mangaData.LastUpdatedAt <= int64(10*60) {
		mangaData.NewAdded++
		return
	}

	var mangaupdatesData models.MangaDetail
	var err error

	sleepTime := 1000
	for {
		mangaupdatesData, err = scrapper.MangaupdatesSeriesDetailByID(mangaData.MangaUpdatesID)
		if err == nil || sleepTime >= 10000 {
			break
		}
		logrus.Errorf("Error updates: %v; %v\n", mangaData.Title, err)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		sleepTime = sleepTime + 1000
	}

	if int(mangaupdatesData.LastChapterInt) > mangaData.MangaLastChapter {
		mangaData.MangaLastChapter = int(mangaupdatesData.LastChapterInt)
		mangaData.NewAdded = 1
		logrus.Infof("Updated: %v - %v", mangaTitle, mangaData.MangaLastChapter)
	} else {
		mangaData.NewAdded++
	}
	mangaData.LastUpdatedAt = time.Now().Unix()

	go rManga.UpdateOneMangaToFireBase(*mangaData)

	mangaDBMutex.Lock()
	mangaDB.MangaDatas[mangaTitle] = mangaData
	mangaDBMutex.Unlock()
}
