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

	var finalMangaTitles []string
	var finalMangaChapters []int

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

		finalMangaTitles = append(finalMangaTitles, sanitizedTitle)
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
			finalMangaChapters = append(finalMangaChapters, latestPage)
		}
	})

	title = strings.Replace(title, " ", "%20", -1)
	c.Visit("https://mangahub.io/search?q=" + title + "&order=POPULAR&genre=all")

	fmt.Println(finalMangaTitles)
	fmt.Println(finalMangaChapters)

	var mangaDB models.MangaDB
	mangaDB.MangaDatas = make(map[string]*models.MangaData)

	for i, title := range finalMangaTitles {
		for j, chapter := range finalMangaChapters {
			if i == j {
				mangaData := models.MangaData{
					MangaLastChapter: chapter,
					AveragePage:      120,
					Status:           "ongoing",
					ImageURL:         "",
					NewAdded:         1,
					Weight:           25000,
					Finder:           "external",
				}
				mangaDB.MangaDatas[title] = &mangaData
			}
		}
	}

	return mangaDB
}
