package main

import (
	"fmt"

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
