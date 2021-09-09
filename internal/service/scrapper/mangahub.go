package scrapper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/lib/app_cache"
	"github.com/umarkotak/go-animapu/internal/models"
)

// SearchMangaTitle from mangahub search page
func SearchMangaTitle(title string) models.MangaDB {
	c := colly.NewCollector()

	var finalMangaTitles map[string]string = make(map[string]string)
	var finalMangaChapters map[string]int = make(map[string]int)
	titleIdx := 1
	chapterIdx := 1

	c.OnHTML(".media-left", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("a img", "src")

		splittedStrings := strings.Split(imageLink, "/")
		lenString := len(splittedStrings)
		var choosenTitle string
		for i, v := range splittedStrings {
			if i == lenString-1 {
				choosenTitle = v
			}
		}
		sanitizedTitle := strings.Replace(choosenTitle, ".jpg", "", -1)
		sanitizedTitle = strings.Replace(sanitizedTitle, ".png", "", -1)

		finalMangaTitles[strconv.Itoa(titleIdx)] = sanitizedTitle
		titleIdx++
	})

	c.OnHTML(".media-body p", func(e *colly.HTMLElement) {
		textDesc := e.Text
		var choosenText string
		var words []string
		var word string
		if strings.Contains(textDesc, "#") {
			choosenText = textDesc
			words = strings.Fields(choosenText)
		}
		for i, v := range words {
			if i == 0 {
				word = v
			}
		}
		sanitizedWord := strings.Replace(word, "#", "", -1)
		latestPageFloat64, _ := strconv.ParseFloat(sanitizedWord, 64)
		latestPage := int(latestPageFloat64)

		if latestPage != 0 {
			finalMangaChapters[strconv.Itoa(chapterIdx)] = latestPage
			chapterIdx++
		}
	})

	title = strings.Replace(title, " ", "%20", -1)
	c.Visit("https://mangahub.io/search?q=" + title + "&order=POPULAR&genre=all")

	// fmt.Println(finalMangaTitles)
	// fmt.Println(finalMangaChapters)

	var mangaDB models.MangaDB
	mangaDB.MangaDatas = make(map[string]*models.MangaData)

	maxWeight := 1000000
	for key, title := range finalMangaTitles {
		weight, _ := strconv.Atoi(key)
		mangaData := models.MangaData{
			MangaLastChapter: finalMangaChapters[key],
			AveragePage:      120,
			Status:           "ongoing",
			ImageURL:         "",
			NewAdded:         1,
			Weight:           maxWeight - weight,
			Finder:           "external",
		}
		mangaDB.MangaDatas[title] = &mangaData
	}

	return mangaDB
}

// GetTodaysMangaTitle get manga titles
func GetTodaysMangaTitle() models.MangaDB {
	appCache := pkgAppCache.GetAppCache()

	res, found := appCache.Get("mangahub_todays_manga")
	if found {
		fmt.Println("FETCH FROM APP CACHE")
		return res.(models.MangaDB)
	}

	var finalMangaTitles map[string]string = make(map[string]string)
	var finalMangaChapters map[string]int = make(map[string]int)
	titleIdx := 1
	chapterIdx := 1

	c := colly.NewCollector()

	c.OnHTML(".panel.panel-default li", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("a img", "src")

		splittedStrings := strings.Split(imageLink, "/")
		lenString := len(splittedStrings)
		var choosenTitle string
		for i, v := range splittedStrings {
			if i == lenString-1 {
				choosenTitle = v
			}
		}
		sanitizedTitle := strings.Replace(choosenTitle, ".jpg", "", -1)
		sanitizedTitle = strings.Replace(sanitizedTitle, ".png", "", -1)

		finalMangaTitles[strconv.Itoa(titleIdx)] = sanitizedTitle
		titleIdx++
	})

	c.OnHTML(".panel.panel-default li a span", func(e *colly.HTMLElement) {
		textDesc := e.Text
		var choosenText string
		var words []string
		var word string
		if strings.Contains(textDesc, "#") {
			choosenText = textDesc
			words = strings.Fields(choosenText)
		}
		for i, v := range words {
			if i == 0 {
				word = v
			}
		}
		sanitizedWord := strings.Replace(word, "#", "", -1)
		latestPage, _ := strconv.Atoi(sanitizedWord)

		if latestPage != 0 {
			finalMangaChapters[strconv.Itoa(chapterIdx)] = latestPage
			chapterIdx++
		}
	})

	c.Visit("https://mangahub.io")

	var mangaDB models.MangaDB
	mangaDB.MangaDatas = make(map[string]*models.MangaData)

	maxWeight := 1000000
	for key, title := range finalMangaTitles {
		weight, _ := strconv.Atoi(key)
		mangaData := models.MangaData{
			MangaLastChapter: finalMangaChapters[key],
			AveragePage:      120,
			Status:           "ongoing",
			ImageURL:         "",
			NewAdded:         1,
			Weight:           maxWeight - weight,
			Finder:           "external",
		}
		mangaDB.MangaDatas[title] = &mangaData
	}

	appCache.Set("mangahub_todays_manga", mangaDB, 5*time.Minute)

	return mangaDB
}

// GetTodaysMangaTitleV2 get manga titles
func GetTodaysMangaTitleV2() models.MangaDB {
	fmt.Println("INCODMING GetTodaysMangaTitleV2")
	appCache := pkgAppCache.GetAppCache()

	res, found := appCache.Get("mangahub_todays_manga")
	if found {
		fmt.Println("FETCH FROM APP CACHE")
		return res.(models.MangaDB)
	}

	var finalMangaTitles map[string]string = make(map[string]string)
	var finalMangaChapters map[string]int = make(map[string]int)
	var rawImageLinks map[string]string = make(map[string]string)
	titleIdx := 1

	c := colly.NewCollector()

	c.OnHTML(".panel.panel-default li", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("a img", "src")

		splittedStrings := strings.Split(imageLink, "/")
		lenString := len(splittedStrings)
		var choosenTitle string
		for i, v := range splittedStrings {
			if i == lenString-1 {
				choosenTitle = v
			}
		}
		sanitizedTitle := strings.Replace(choosenTitle, ".jpg", "", -1)
		sanitizedTitle = strings.Replace(sanitizedTitle, ".png", "", -1)
		sanitizedTitle = strings.Replace(sanitizedTitle, ".jpeg", "", -1)

		textDesc := e.ChildText("li a span")
		var choosenText string
		var words []string
		var word string
		if strings.Contains(textDesc, "#") {
			choosenText = textDesc
			words = strings.Fields(choosenText)
		}
		for i, v := range words {
			if i == 0 {
				word = v
			}
		}
		sanitizedWord := strings.Replace(word, "#", "", -1)
		latestPageFloat64, _ := strconv.ParseFloat(sanitizedWord, 64)
		latestPage := int(latestPageFloat64)
		finalMangaChapters[strconv.Itoa(titleIdx)] = latestPage

		finalMangaTitles[strconv.Itoa(titleIdx)] = sanitizedTitle
		rawImageLinks[strconv.Itoa(titleIdx)] = imageLink

		// fmt.Printf("%v | %v | %v \n", titleIdx, sanitizedTitle, latestPage)

		titleIdx++
	})

	c.Visit("https://mangahub.io")

	var mangaDB models.MangaDB
	mangaDB.MangaDatas = make(map[string]*models.MangaData)

	maxWeight := 1000000
	for key, title := range finalMangaTitles {
		weight, _ := strconv.Atoi(key)
		mangaData := models.MangaData{
			MangaLastChapter: finalMangaChapters[key],
			AveragePage:      120,
			Status:           "ongoing",
			ImageURL:         rawImageLinks[key],
			NewAdded:         1,
			Weight:           maxWeight - weight,
			Finder:           "external",
		}
		mangaDB.MangaDatas[title] = &mangaData
	}

	appCache.Set("mangahub_todays_manga", mangaDB, 5*time.Minute)
	_, b := appCache.Get("mangahub_todays_manga")
	fmt.Println("APP CACHE IS", b)

	return mangaDB
}

// GetMangaDetailV1 get manga detail
func GetMangaDetailV1(title string) models.MangaDetail {
	mangaDetail := models.MangaDetail{}

	url := fmt.Sprintf("https://mangahub.io/manga/%v", title)

	c := colly.NewCollector()

	c.OnHTML("ul.list-group li", func(e *colly.HTMLElement) {
		rawChapter := e.ChildAttr("a", "href")
		candidates := strings.Split(rawChapter, "-")
		chapter := candidates[len(candidates)-1]

		mangaDetail.Chapters = append(mangaDetail.Chapters, chapter)
	})

	tempGenres := []string{}
	c.OnHTML(".label.genre-label", func(e *colly.HTMLElement) {
		tempGenres = append(tempGenres, e.Text)
	})

	c.OnHTML(".img-responsive.manga-thumb", func(e *colly.HTMLElement) {
		mangaDetail.ImageURL = e.Attr("src")
	})

	c.Visit(url)

	mangaDetail.Title = title
	mangaDetail.Genres += strings.Join(tempGenres, ", ")

	return mangaDetail
}
