package scrapper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/umarkotak/go-animapu/internal/models"
)

var (
	maidMyHost = "https://www.maid.my.id"
)

func ScrapMaidMyHomePage() models.MangaDB {
	c := colly.NewCollector()

	mangaDatas := map[string]*models.MangaData{}

	c.OnHTML(".flexbox3-content", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("img", "src")
		mangaLink := e.ChildAttr("a", "href")
		splittedLink := strings.Split(mangaLink, "/manga/")
		mangaTitle := splittedLink[1]
		mangaTitle = strings.ReplaceAll(mangaTitle, "/", "")

		tempMangaData := models.MangaData{
			CompactTitle:     mangaTitle,
			MangaLastChapter: 150,
			AveragePage:      150,
			ImageURL:         imageLink,
		}
		mangaDatas[mangaTitle] = &tempMangaData
	})

	c.Visit(fmt.Sprintf("%v/page/1/", maidMyHost))

	c.Visit(fmt.Sprintf("%v/page/2/", maidMyHost))

	c.Visit(fmt.Sprintf("%v/page/3/", maidMyHost))

	c.Visit(fmt.Sprintf("%v/page/4/", maidMyHost))

	c.Visit(fmt.Sprintf("%v/page/5/", maidMyHost))

	mangaDB := models.MangaDB{
		MangaDatas: mangaDatas,
	}
	return mangaDB
}

func ScrapMaidMyMangaSearchPage(query string) models.MangaDB {
	mangaDB := models.MangaDB{}
	return mangaDB
}

func ScrapMaidMyMangaDetailPage(title string) models.MangaDetail {
	c := colly.NewCollector()

	mangaDetail := models.MangaDetail{
		Title: title,
	}

	c.OnHTML(".flexch-infoz", func(e *colly.HTMLElement) {
		tempChapterLink := e.ChildAttr("a", "href")
		tempChapterLinkArr := strings.Split(tempChapterLink, "www.maid.my.id")
		chapterLink := tempChapterLinkArr[1]
		chapterLink = strings.ReplaceAll(chapterLink, "/", "")

		tempChapterNo := e.ChildText("span.ch")
		tempChapterNoArr := strings.Split(tempChapterNo, " ")
		chapterNo := tempChapterNoArr[1]

		mangaDetail.ChapterLinks = append(mangaDetail.ChapterLinks, chapterLink)
		mangaDetail.Chapters = append(mangaDetail.Chapters, chapterNo)
	})

	c.OnHTML(".series-thumb", func(e *colly.HTMLElement) {
		imageUrl := e.ChildAttr("img", "src")

		mangaDetail.ImageURL = imageUrl
	})

	c.OnHTML(".series-synops", func(e *colly.HTMLElement) {
		description := e.ChildText("p")

		mangaDetail.Description = description
	})

	c.OnHTML(".series-genres", func(e *colly.HTMLElement) {
		tempGenres := e.ChildAttrs("a", "href")
		tempGenresClean := []string{}

		for _, tempGenre := range tempGenres {
			tempGenreArr := strings.Split(tempGenre, "genres/")
			tempGenreClean := tempGenreArr[1]
			tempGenreClean = strings.ReplaceAll(tempGenreClean, "/", "")
			tempGenresClean = append(tempGenresClean, tempGenreClean)
		}

		genres := strings.Join(tempGenresClean, ", ")

		mangaDetail.Genres = genres
	})

	c.Visit(fmt.Sprintf("%v/manga/%v/", maidMyHost, title))

	return mangaDetail
}

func ScrapMaidMyMangaChapterDetailPage(title, chapter string) models.MangaChapterDetail {
	return models.MangaChapterDetail{}
}