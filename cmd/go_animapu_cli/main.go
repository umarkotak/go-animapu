package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/umarkotak/go-animapu/internal/models"
	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
)

func main() {
	fmt.Println("Welcome to go-animapu CLI")

	initBaseConfiguration()

	var mangaDB models.MangaDB

	if true == false {
		mangaDB = rManga.GetMangaFromJSON()
		mangaDB = sManga.UpdateMangaChapters(mangaDB)
		mangaDB = rManga.UpdateMangaToFireBase(mangaDB)
	}

	mangaDB = rManga.GetMangaFromFireBaseV2()
	mangaDB = rManga.UpdateMangaToFireBase(mangaDB)

}

func initBaseConfiguration() {
	godotenv.Load(".env")
}
