package scrapper

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/umarkotak/go-animapu/internal/models"
)

func ScrapKlikMangaHomePage() models.MangaDB {
	c := colly.NewCollector()

	mangaDataKeys := []string{}
	mangaDatas := map[string]*models.MangaData{}

	weight := 10000

	c.OnHTML("#loop-content > div > div > div > div", func(e *colly.HTMLElement) {
		compactTitle := e.ChildText("div.item-summary > div.post-title.font-title > h3 > a")

		mangaLink := e.ChildAttr("div.item-summary > div.post-title.font-title > h3 > a", "href")
		mangaTitle := strings.Replace(mangaLink, "https://klikmanga.com/manga/", "", -1)
		mangaTitle = strings.Replace(mangaTitle, "/", "", -1)

		chapter := e.ChildText("div.item-summary > div.list-chapter > div > span.chapter.font-meta")
		chapter = strings.Replace(chapter, "Chapter ", "", -1)
		chapterBreaks := strings.Split(chapter, " ")
		chapter = chapterBreaks[0]
		chapterFloat, _ := strconv.ParseFloat(chapter, 64)

		imageURL := e.ChildAttr("div.item-thumb.hover-details.c-image-hover > a > img", "src")

		mangaData := models.MangaData{
			Title:            mangaTitle,
			CompactTitle:     compactTitle,
			MangaLastChapter: int(chapterFloat),
			Weight:           weight,
			ImageURL:         imageURL,
		}
		mangaDatas[mangaTitle] = &mangaData
		mangaDataKeys = append(mangaDataKeys, mangaTitle)
		weight--
	})

	c.Visit("https://klikmanga.com/")

	mangaDB := models.MangaDB{
		MangaDataKeys: mangaDataKeys,
		MangaDatas:    mangaDatas,
	}
	return mangaDB
}
