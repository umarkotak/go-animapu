package scrapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
)

func ScrapMangabatList(page int64) []models.MangaData {
	c := colly.NewCollector()

	mangaDatas := []models.MangaData{}
	weight := 10000

	c.OnHTML("body div.body-site div.container.container-main div.container-main-left div.panel-list-story .list-story-item", func(e *colly.HTMLElement) {
		mangaTitle := e.ChildText("div > h3 > a")

		compactTitle := strings.Replace(e.ChildAttr("div > h3 > a", "href"), "https://read.mangabat.com/", "", -1)
		compactTitle = strings.Replace(compactTitle, "https://m.mangabat.com/", "", -1)

		imageURL := e.ChildAttr("a > img", "src")

		lastChapterText := e.ChildText("div > a:nth-child(2)")

		lastChapterID := strings.Replace(e.ChildAttr("div > a:nth-child(2)", "href"), "https://read.mangabat.com/", "", -1)
		lastChapterID = strings.Replace(lastChapterID, "https://m.mangabat.com/", "", -1)

		chapterVals := strings.Split(lastChapterText, ":")
		chapterValString := strings.Replace(chapterVals[0], "Chapter ", "", -1)
		chapterVal, _ := strconv.ParseInt(chapterValString, 10, 64)

		mangaData := models.MangaData{
			Title:            mangaTitle,
			CompactTitle:     compactTitle,
			MangaLastChapter: int(chapterVal),
			Weight:           weight,
			ImageURL:         imageURL,
			LastChapterID:    lastChapterID,
			LastChapterText:  lastChapterText,
		}
		mangaDatas = append(mangaDatas, mangaData)
		weight--
	})

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", page))
	if err != nil {
		logrus.Errorf("ScrapMangabatHomePage: %v\n", err)
	}

	return mangaDatas
}

func ScrapMangabatDetail(mangaID string) models.MangaDetail {
	c := colly.NewCollector()

	mangaDetail := models.MangaDetail{
		Chapters:     []string{},
		ChaptersInt:  []int64{},
		ChapterLinks: []string{},
		ChapterObjs:  []models.ChapterObj{},
	}

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info > div.story-info-right > h1", func(e *colly.HTMLElement) {
		mangaDetail.Title = e.Text
	})

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("a", "href")
		chapterLink = strings.Replace(chapterLink, "https://read.mangabat.com/", "", -1)
		chapterLink = strings.Replace(chapterLink, "https://m.mangabat.com/", "", -1)

		chapterObj := models.ChapterObj{
			Title: e.ChildText("a"),
			Link:  chapterLink,
		}

		mangaDetail.ChapterObjs = append(mangaDetail.ChapterObjs, chapterObj)
	})

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/%v", mangaID))
	if err != nil {
		logrus.Errorf("ScrapMangabatDetail: %v\n", err)
	}

	if mangaDetail.Title == "" {
		err := c.Visit(fmt.Sprintf("https://read.mangabat.com/%v", mangaID))
		if err != nil {
			logrus.Errorf("ScrapMangabatDetail: %v\n", err)
		}
	}

	return mangaDetail
}

func ScrapMangabatChapterDetail(mangaID, chapterID string) models.MangaChapterDetail {
	c := colly.NewCollector()

	mangaChapterDetail := models.MangaChapterDetail{
		ChapterObjs: []models.ChapterObj{},
		Images:      []string{},
	}

	c.OnHTML("body > div.body-site > div.container-chapter-reader > img", func(e *colly.HTMLElement) {
		mangaChapterDetail.Images = append(mangaChapterDetail.Images, fmt.Sprintf("https://go-animapu.herokuapp.com/image_proxy/%v", e.Attr("src")))
	})

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/%v", chapterID))
	if err != nil {
		logrus.Errorf("ScrapMangabatDetail: %v\n", err)
	}

	return mangaChapterDetail
}
