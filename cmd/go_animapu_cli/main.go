package main

import (
	"fmt"
	"strconv"
	"strings"

	colly "github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/umarkotak/go-animapu/internal/models"
	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	rUser "github.com/umarkotak/go-animapu/internal/repository/user"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
)

var mangaDB models.MangaDB

func main() {
	fmt.Println("Welcome to go-animapu CLI")

	initBaseConfiguration()
	c := colly.NewCollector()

	var finalMangaTitles []string
	var finalMangaPages []int

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
			finalMangaPages = append(finalMangaPages, latestPage)
		}
	})

	c.Visit("https://mangahub.io/search?q=blea&order=POPULAR&genre=all")

	fmt.Println(finalMangaTitles)
	fmt.Println(finalMangaPages)

	fmt.Println("Thanks for using go-animapu CLI")
}

func initBaseConfiguration() {
	godotenv.Load(".env")
}

func learnMangaJSON() {
	mangaDB = rManga.GetMangaFromJSON()
	mangaDB = sManga.UpdateMangaChapters(mangaDB)
	mangaDB = rManga.UpdateMangaToFireBase(mangaDB)
}

func learnMangaFirebase() {
	mangaDB = rManga.GetMangaFromFireBaseV2()
	mangaDB = rManga.UpdateMangaToFireBase(mangaDB)
}

func learnUserFirebase() {
	userData := models.UserData{
		Username: "hello",
		Password: "goodbye",
		ReadHistories: map[string]*models.ReadHistory{
			"a-returner-s-magic-should-be-special": {
				MangaTitle:    "a-returner-s-magic-should-be-special",
				LastChapter:   10,
				LastReadTime:  "2020-01-01T00:00",
				LastReadTimeI: 1000,
			},
		},
	}

	userData = rUser.SetUserToFirebase(userData)
	fUser := rUser.GetUserByUsernameFromFirebase("hello")
	fmt.Println(fUser)
}
