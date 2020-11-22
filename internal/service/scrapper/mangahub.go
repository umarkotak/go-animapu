package scrapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
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
		latestPage, _ := strconv.Atoi(sanitizedWord)

		if latestPage != 0 {
			finalMangaChapters[strconv.Itoa(chapterIdx)] = latestPage
			chapterIdx++
		}
	})

	title = strings.Replace(title, " ", "%20", -1)
	c.Visit("https://mangahub.io/search?q=" + title + "&order=POPULAR&genre=all")
	// c.Visit("https://mangahub.io")

	fmt.Println(finalMangaTitles)
	fmt.Println(finalMangaChapters)

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

	return mangaDB
}
