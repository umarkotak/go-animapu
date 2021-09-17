package scrapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
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

	err := c.Visit("https://klikmanga.com/")
	if err != nil {
		logrus.Errorf("ScrapKlikMangaHomePage: %v\n", err)
	}

	mangaDB := models.MangaDB{
		MangaDataKeys: mangaDataKeys,
		MangaDatas:    mangaDatas,
	}
	return mangaDB
}

func ScrapKlikMangaDetailPage(title string) models.MangaDetail {
	c := colly.NewCollector()
	url := fmt.Sprintf("https://klikmanga.com/manga/%v", title)
	maxChapter := 0

	mangaDetail := models.MangaDetail{
		Title:        title,
		ChapterLinks: []string{},
		Chapters:     []string{},
		ChaptersInt:  []int64{},
	}

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.post-title > h1", func(e *colly.HTMLElement) {
		mangaDetail.CompactTitle = e.Text
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.c-page-content.style-1 > div > div > div > div.main-col.col-md-8.col-sm-8 > div > div.c-page > div > div.page-content-listing.single-page > div > ul > li", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("a", "href")

		chapterLinkId := strings.Replace(chapterLink, url, "", -1)
		chapterLinkId = strings.Replace(chapterLinkId, "/", "", -1)

		chapterString := e.ChildText("a")
		chapterString = strings.Replace(chapterString, "Chapter ", "", -1)
		chapterInt, _ := strconv.ParseFloat(chapterString, 64)

		mangaDetail.ChapterLinks = append(mangaDetail.ChapterLinks, chapterLink)
		mangaDetail.Chapters = append(mangaDetail.Chapters, chapterLinkId)
		mangaDetail.ChaptersInt = append(mangaDetail.ChaptersInt, int64(chapterInt))

		if int(chapterInt) > maxChapter {
			mangaDetail.LastChapter = chapterLinkId
			mangaDetail.LastChapterInt = int64(chapterInt)
			maxChapter = int(chapterInt)
		}
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.tab-summary > div.summary_image > a > img", func(e *colly.HTMLElement) {
		mangaDetail.ImageURL = e.Attr("src")
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.c-page-content.style-1 > div > div > div > div.main-col.col-md-8.col-sm-8 > div > div.c-page > div > div.description-summary > div > p", func(e *colly.HTMLElement) {
		mangaDetail.Description = e.Text
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.tab-summary > div.summary_content_wrap > div > div.post-content > div:nth-child(8) > div.summary-content > div", func(e *colly.HTMLElement) {
		mangaDetail.Genres = e.Text
	})

	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("ScrapKlikMangaDetailPage: %v\n", err)
	}

	return mangaDetail
}

func ScrapKlikMangaChapterDetailPage(title, chapter string) models.MangaChapterDetail {
	c := colly.NewCollector()
	url := fmt.Sprintf("https://klikmanga.com/manga/%v/%v/?style=list", title, chapter)

	mangaChapter := models.MangaChapterDetail{
		Title:  title,
		Images: []string{},
	}

	c.OnHTML("body > div.wrap > div > div.site-content > div > div > div > div > div > div > div.c-blog-post > div.entry-content > div > div > div.reading-content > div", func(e *colly.HTMLElement) {
		image := e.ChildAttr("img.wp-manga-chapter-img", "src")
		mangaChapter.Images = append(mangaChapter.Images, image)
	})

	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("ScrapKlikMangaDetailPage: %v\n", err)
	}

	return mangaChapter
}

func ScrapKlikMangaSearch(searchParams models.KlikMangaSearchParams) models.MangaDB {
	c := colly.NewCollector()
	url := fmt.Sprintf(
		"https://klikmanga.com/?s=%v&post_type=wp-manga&op=&author=&artist=&release=&adult=&genre%5B%5D=%v&status%5B%5D=%v",
		searchParams.Title, searchParams.Genre, searchParams.Status,
	)

	mangaDataKeys := []string{}
	mangaDatas := map[string]*models.MangaData{}
	weight := 10000

	c.OnHTML("body > div.wrap > div > div.site-content > div.c-page-content > div > div > div > div > div.main-col-inner > div > div.tab-content-wrap > div > div", func(e *colly.HTMLElement) {
		mangaLink := e.ChildAttr("div div a", "href")

		mangaTitle := strings.Replace(mangaLink, "https://klikmanga.com/manga/", "", -1)
		mangaTitle = strings.Replace(mangaTitle, "/", "", -1)

		compactTitle := e.ChildText("div.col-8.col-12.col-md-10 > div.tab-summary > div.post-title > h3 > a")

		lastChapterLink := e.ChildAttr("div.col-8.col-12.col-md-10 > div.tab-meta > div.meta-item.latest-chap > span.font-meta.chapter > a", "href")
		prefix := fmt.Sprintf("https://klikmanga.com/manga/%v", mangaTitle)
		lastChapterID := strings.Replace(lastChapterLink, prefix, "", -1)
		lastChapterID = strings.Replace(lastChapterID, "/", "", -1)

		lastChapterInt, _ := strconv.ParseFloat(lastChapterID, 64)

		imageURL := e.ChildAttr("div.col-4.col-12.col-md-2 > div > a > img", "src")

		mangaData := models.MangaData{
			Title:            mangaTitle,
			ImageURL:         imageURL,
			CompactTitle:     compactTitle,
			LastChapterID:    lastChapterID,
			MangaLastChapter: int(lastChapterInt),
			Weight:           weight,
		}

		mangaDatas[mangaData.Title] = &mangaData
		mangaDataKeys = append(mangaDataKeys, mangaData.Title)
		weight--
	})

	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("ScrapKlikMangaSearch: %v\n", err)
	}

	mangaDB := models.MangaDB{
		MangaDataKeys: mangaDataKeys,
		MangaDatas:    mangaDatas,
	}
	return mangaDB
}
