package main

import (
	"fmt"

	rManga "github.com/umarkotak/go-animapu/internal/repository/manga"
	sManga "github.com/umarkotak/go-animapu/internal/service/manga"
)

func main() {
	fmt.Println("Welcome to go-animapu CLI")

	mangaDB := rManga.GetMangaFromJSON()
	mangaDB = sManga.UpdateMangaChapters(mangaDB)
	mangaDB = rManga.UpdateMangaToJSON(mangaDB)

}
