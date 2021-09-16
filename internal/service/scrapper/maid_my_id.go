package scrapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
)

var (
	maidMyHost = "https://www.maid.my.id"
)

func ScrapMaidMyHomePage() models.MangaDB {
	c := colly.NewCollector()

	mangaDataKeys := []string{}
	mangaDatas := map[string]*models.MangaData{}

	maxWeight := 1000000

	c.OnHTML("body > main > div > div.container > div.flexbox4 > div", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("div.flexbox4-content a div.flexbox4-thumb img", "src")
		mangaCompactTitle := e.ChildAttr("div.flexbox4-content a div.flexbox4-thumb img", "title")

		mangaLink := e.ChildAttr("div.flexbox4-content a", "href")
		splittedLink := strings.Split(mangaLink, "/manga/")
		mangaTitle := splittedLink[1]
		mangaTitle = strings.ReplaceAll(mangaTitle, "/", "")

		mangaLastChapterRaw := e.ChildText("div > div > ul > li:nth-child(1) > a")
		mangaLastChapterRaw = strings.ReplaceAll(mangaLastChapterRaw, "Ch.", "")
		mangaLastChapterRaw = strings.ReplaceAll(mangaLastChapterRaw, " ", "")
		mangaLastChapter, err := strconv.ParseFloat(mangaLastChapterRaw, 64)
		if err != nil {
			logrus.Error(mangaLastChapterRaw, err)
			mangaLastChapter = 0
		}

		tempMangaData := models.MangaData{
			Title:            mangaTitle,
			CompactTitle:     mangaCompactTitle,
			MangaLastChapter: int(mangaLastChapter),
			AveragePage:      100,
			ImageURL:         imageLink,
			Status:           "ongoing",
			Weight:           maxWeight,
		}
		fmt.Printf("%+v\n", tempMangaData)

		_, ok := mangaDatas[mangaTitle]
		if !ok {
			mangaDatas[mangaTitle] = &tempMangaData
			mangaDataKeys = append(mangaDataKeys, mangaTitle)
			maxWeight--
		}
	})

	c.Visit(fmt.Sprintf("%v/page/1/", maidMyHost))
	c.Visit(fmt.Sprintf("%v/page/2/", maidMyHost))
	c.Visit(fmt.Sprintf("%v/page/3/", maidMyHost))
	c.Visit(fmt.Sprintf("%v/page/4/", maidMyHost))
	c.Visit(fmt.Sprintf("%v/page/5/", maidMyHost))

	mangaDB := models.MangaDB{
		MangaDataKeys: mangaDataKeys,
		MangaDatas:    mangaDatas,
	}
	return mangaDB
}

func ScrapMaidMyMangaSearchPage(query string) models.MangaDB {
	c := colly.NewCollector()

	mangaDatas := map[string]*models.MangaData{}

	maxWeight := 1000000
	idx := 0

	c.OnHTML(".flexbox2-content", func(e *colly.HTMLElement) {
		imageLink := e.ChildAttr("img", "src")
		mangaLink := e.ChildAttr("a", "href")
		splittedLink := strings.Split(mangaLink, "/manga/")
		mangaTitle := splittedLink[1]
		mangaTitle = strings.ReplaceAll(mangaTitle, "/", "")
		lastChapterS := e.ChildText(".season")
		lastChapterSArr := strings.Split(lastChapterS, " ")
		lastChapter := int64(150)
		if len(lastChapterSArr) == 2 {
			lastChapter, _ = strconv.ParseInt(lastChapterSArr[1], 10, 64)
		}

		weight := maxWeight - idx
		idx += 1

		tempMangaData := models.MangaData{
			CompactTitle:     mangaTitle,
			MangaLastChapter: int(lastChapter),
			AveragePage:      150,
			ImageURL:         imageLink,
			Status:           "ongoing",
			Weight:           weight,
		}
		mangaDatas[mangaTitle] = &tempMangaData
	})

	c.Visit(fmt.Sprintf("%v/?s=%v", maidMyHost, query))

	mangaDB := models.MangaDB{
		MangaDatas: mangaDatas,
	}
	return mangaDB
}

func ScrapMaidMyMangaDetailPage(title string) models.MangaDetail {
	c := colly.NewCollector()

	mangaDetail := models.MangaDetail{
		Title: title,
	}

	maxChapter := int64(0)
	c.OnHTML(".flexch-infoz", func(e *colly.HTMLElement) {
		tempChapterLink := e.ChildAttr("a", "href")
		tempChapterLinkArr := strings.Split(tempChapterLink, "www.maid.my.id")
		chapterLink := tempChapterLinkArr[1]
		chapterLink = strings.ReplaceAll(chapterLink, "/", "")

		tempChapterNo := e.ChildText("span")
		tempSuffix := e.ChildText("span span")
		tempChapterNo = strings.ReplaceAll(tempChapterNo, tempSuffix, "")
		tempChapterNos := strings.Split(tempChapterNo, " ")
		lastChapterNoString := "0"
		if len(tempChapterNos) == 2 {
			lastChapterNoString = tempChapterNos[1]
		}

		lastChapterNoInt, _ := strconv.ParseFloat(lastChapterNoString, 64)
		if lastChapterNoInt > float64(maxChapter) {
			maxChapter = int64(lastChapterNoInt)
			mangaDetail.LastChapter = fmt.Sprintf("%v", maxChapter)
		}

		mangaDetail.ChapterLinks = append(mangaDetail.ChapterLinks, chapterLink)
		mangaDetail.Chapters = append(mangaDetail.Chapters, lastChapterNoString)
		mangaDetail.ChaptersInt = append(mangaDetail.ChaptersInt, int64(lastChapterNoInt))
		mangaDetail.LastChapterInt = maxChapter
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

	c.OnHTML("body > main > div > div > div.container > div > div.series-flexright > div.series-title > h2", func(e *colly.HTMLElement) {
		mangaDetail.CompactTitle = e.Text
	})

	c.Visit(fmt.Sprintf("%v/manga/%v/", maidMyHost, title))

	return mangaDetail
}

func ScrapMaidMyMangaChapterDetailPage(title, chapter string) models.MangaChapterDetail {
	c := colly.NewCollector()

	mangaChapterDetail := models.MangaChapterDetail{
		Title: title,
	}

	fmt.Println(fmt.Sprintf("%v/%v/", maidMyHost, chapter))

	c.OnHTML(".reader-area", func(e *colly.HTMLElement) {
		imageUrls := e.ChildAttrs("img", "src")

		mangaChapterDetail.Images = imageUrls
	})

	c.Visit(fmt.Sprintf("%v/%v/", maidMyHost, chapter))
	return mangaChapterDetail
}
