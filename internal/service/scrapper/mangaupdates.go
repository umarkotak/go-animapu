package scrapper

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
	"go4.org/sort"
)

type (
	MangaUpdatesResult struct {
		Title string
		Raw   map[string]interface{}
	}
)

func MangaupdatesGetSeries(mangaID string) ([]MangaUpdatesResult, error) {
	url := fmt.Sprintf("https://www.mangaupdates.com/series.html?id=%v", mangaID)

	mangaUpdatesResults := []MangaUpdatesResult{}

	c := colly.NewCollector()

	c.OnHTML("span.releasestitle.tabletitle", func(e *colly.HTMLElement) {
		tempMangaResult := MangaUpdatesResult{}
		tempRaw := map[string]interface{}{}
		tempRaw["title"] = e.Text
		mangaUpdatesResults = append(mangaUpdatesResults, tempMangaResult)
	})

	err := c.Visit(url)

	return mangaUpdatesResults, err
}

func MangaupdatesGetReleases() (models.MangaDB, error) {
	mangaDatas := map[string]*models.MangaData{}
	mangaDB := models.MangaDB{MangaDatas: mangaDatas}

	c := colly.NewCollector()

	prevTitle := ""
	perveMangadexID := ""
	idx := 100000

	c.OnHTML("div.alt.p-1 div.row.no-gutters div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "col-6 pbreak" {
			prevTitle = e.ChildText("a")
			perveMangadexID = e.ChildAttr("a", "href")
			mangaupdatesIDRaw := strings.Split(perveMangadexID, "series.html?id=")
			if len(mangaupdatesIDRaw) >= 2 {
				perveMangadexID = mangaupdatesIDRaw[1]
			}
		}

		if e.Attr("class") == "col-2 pl-1 pbreak" {
			mangahubTitle := mangahubTitleCostructor(prevTitle)
			mangaLastChapters := strings.Split(e.Text, "c.")
			var mangaLastChapterString string
			if len(mangaLastChapters) > 0 {
				mangaLastChapterString = mangaLastChapters[len(mangaLastChapters)-1]
			} else {
				mangaLastChapterString = ""
			}

			mangaLastChapterBreaked := strings.Split(mangaLastChapterString, "-")
			if len(mangaLastChapterBreaked) > 0 {
				mangaLastChapterString = mangaLastChapterBreaked[len(mangaLastChapterBreaked)-1]
			}

			mangaLastChapterString = strings.Replace(mangaLastChapterString, "a", "", -1)
			mangaLastChapterString = strings.Replace(mangaLastChapterString, "b", "", -1)
			mangaLastChapterString = strings.Replace(mangaLastChapterString, "c", "", -1)
			var mangaLastChapter float64
			mangaLastChapter, err := strconv.ParseFloat(mangaLastChapterString, 64)
			if err != nil {
				mangaLastChapter = 150
			}

			mangaData := models.MangaData{
				Title:            mangahubTitle,
				CompactTitle:     prevTitle,
				MangaLastChapter: int(mangaLastChapter),
				Weight:           idx,
				MangaUpdatesID:   perveMangadexID,
			}
			_, found := mangaDatas[mangahubTitle]
			if !found {
				mangaDatas[mangahubTitle] = &mangaData
			}
			idx--
		}

	})

	c.SetRequestTimeout(60 * time.Second)
	var err error
	err = c.Visit("https://www.mangaupdates.com/releases.html?page=1")
	if err != nil {
		logrus.Errorf("There is some error: %v", err)
		return mangaDB, err
	}
	err = c.Visit("https://www.mangaupdates.com/releases.html?page=2")
	if err != nil {
		logrus.Errorf("There is some error: %v", err)
		return mangaDB, err
	}

	return mangaDB, nil
}

func MangaupdatesSeriesDetailByID(mangaupdateID string) (models.MangaDetail, error) {
	mangaDetail := models.MangaDetail{}

	url := fmt.Sprintf("https://www.mangaupdates.com/releases.html?stype=series&search=%v&page=1&perpage=100&orderby=date&asc=desc", mangaupdateID)

	c := colly.NewCollector()

	c.OnHTML("div#main_content div.p-2.pt-2.pb-2.text div.row.no-gutters div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "col-4 text" {
			if mangaDetail.CompactTitle == "" {
				mangaDetail.CompactTitle = e.ChildText("a")
				mangaDetail.DetailLink = e.ChildAttr("a", "href")
				mangaDetail.Title = mangahubTitleCostructor(mangaDetail.CompactTitle)

				mangaupdatesIDRaw := strings.Split(mangaDetail.DetailLink, "series.html?id=")
				if len(mangaupdatesIDRaw) >= 2 {
					mangaDetail.MangaUpdatesID = mangaupdatesIDRaw[1]
				}
			}
		}
		if e.Attr("class") == "col-1 text text-center" {
			chapter := e.ChildText("span")

			chapterSplitted := strings.Split(chapter, "-")

			if len(chapterSplitted) > 1 {
				chapter = chapterSplitted[len(chapterSplitted)-1]
			}

			if chapter != "" {
				chapterInt, err := strconv.ParseInt(chapter, 10, 64)
				if err == nil {
					if chapterInt > mangaDetail.LastChapterInt {
						mangaDetail.LastChapterInt = chapterInt
						mangaDetail.LastChapter = fmt.Sprintf("%v", chapterInt)
					}

					mangaDetail.Chapters = append(mangaDetail.Chapters, fmt.Sprintf("%v", chapterInt))
					mangaDetail.ChaptersInt = append(mangaDetail.ChaptersInt, chapterInt)
				}
			}
		}
	})

	c.SetRequestTimeout(60 * time.Second)
	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("There is some error: %v; %v", mangaupdateID, err)
		return mangaDetail, err
	}

	if mangaDetail.DetailLink != "" {
		cDetail := colly.NewCollector()

		cDetail.OnHTML("center > img", func(e *colly.HTMLElement) {
			mangaDetail.ImageURL = e.Attr("src")
		})

		cDetail.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(3) > div:nth-child(2)", func(e *colly.HTMLElement) {
			mangaDetail.Description = e.Text
		})

		cDetail.Visit(mangaDetail.DetailLink)
	}

	sort.Slice(mangaDetail.ChaptersInt, func(i, j int) bool { return mangaDetail.ChaptersInt[i] < mangaDetail.ChaptersInt[j] })

	return mangaDetail, nil
}

func MangaupdatesSearch(title string) (models.MangaDB, error) {
	mangaDatas := map[string]*models.MangaData{}
	mangaDB := models.MangaDB{MangaDatas: mangaDatas, MangaDataKeys: []string{}}

	url := fmt.Sprintf("https://www.mangaupdates.com/series.html?search=%v", url.QueryEscape(title))

	c := colly.NewCollector()

	idx := 100000

	c.OnHTML("#main_content > div.p-2.pt-2.pb-2.text > div:nth-child(2) > div", func(e *colly.HTMLElement) {
		mangaData := models.MangaData{}

		title := e.ChildText("div > div.col.text.p-1.pl-3 > div > div:nth-child(1) > a > u > b")

		if title == "" {
			return
		}

		mangaData.CompactTitle = title
		mangaData.Title = mangahubTitleCostructor(title)

		mangaupdatesSeriesLink := e.ChildAttr("div.col-auto.align-self-center.series_thumb.p-0 > a", "href")
		mangaupdatesIDRaw := strings.Split(mangaupdatesSeriesLink, "series.html?id=")

		mangaImageUrl := e.ChildAttr("div.col-auto.align-self-center.series_thumb.p-0 > a > img", "src")

		if len(mangaupdatesIDRaw) >= 2 {
			mangaData.MangaUpdatesID = mangaupdatesIDRaw[1]
			mangaData.ImageURL = mangaImageUrl
		}
		mangaData.Weight = idx
		idx--

		mangaDatas[mangaData.Title] = &mangaData
		mangaDB.MangaDataKeys = append(mangaDB.MangaDataKeys, mangaData.Title)
	})

	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("There is some error: %v; %v", title, err)
		return mangaDB, err
	}

	return mangaDB, nil
}

func MangaupdatesSeriesDetailByTitle(title string) (models.MangaDetail, error) {
	mangaDetail := models.MangaDetail{}

	url := fmt.Sprintf("https://www.mangaupdates.com/series.html?search=%v", title)
	c := colly.NewCollector()

	c.OnHTML("#main_content > div.p-2.pt-2.pb-2.text > div:nth-child(2) > div:nth-child(3) > div > div.col.text.p-1.pl-3 > div > div:nth-child(1) > a", func(e *colly.HTMLElement) {
		mangaDetail.DetailLink = e.Attr("href")
		mangaupdatesIDRaw := strings.Split(mangaDetail.DetailLink, "series.html?id=")
		if len(mangaupdatesIDRaw) > 0 {
			mangaDetail.MangaUpdatesID = mangaupdatesIDRaw[len(mangaupdatesIDRaw)-1]
		}
	})

	c.Visit(url)

	cDetail := colly.NewCollector()

	cDetail.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div.col-12.p-2 > span.releasestitle.tabletitle", func(e *colly.HTMLElement) {
		mangaDetail.CompactTitle = e.Text
		mangaDetail.Title = mangahubTitleCostructor(mangaDetail.CompactTitle)
	})

	cDetail.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(4) > div:nth-child(5)", func(e *colly.HTMLElement) {
		mangaDetail.Genres = e.Text
	})

	cDetail.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(4) > div:nth-child(2) > center > img", func(e *colly.HTMLElement) {
		mangaDetail.ImageURL = e.Attr("href")
	})

	cDetail.OnHTML("#div_desc_link", func(e *colly.HTMLElement) {
		mangaDetail.Description = e.Text
	})

	cDetail.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(3) > div:nth-child(17) > i:nth-child(1)", func(e *colly.HTMLElement) {
		chapter := e.Text
		chapterSplitted := strings.Split(chapter, "-")
		if len(chapterSplitted) > 1 {
			chapter = chapterSplitted[len(chapterSplitted)-1]
		}

		mangaDetail.LastChapter = chapter

		chapterFloat, _ := strconv.ParseFloat(mangaDetail.LastChapter, 64)

		mangaDetail.LastChapterInt = int64(chapterFloat)

		mangaDetail.ChaptersInt = []int64{}
		for i := int64(1); i <= mangaDetail.LastChapterInt; i++ {
			mangaDetail.ChaptersInt = append(mangaDetail.ChaptersInt, i)
		}
	})

	cDetail.Visit(mangaDetail.DetailLink)

	return mangaDetail, nil
}

func MangaupdatesGetReleasesV2() ([]models.MangaData, error) {
	animapuMangas := []models.MangaData{}

	c := colly.NewCollector()

	prevTitle := ""
	perveMangadexID := ""
	idx := 100000

	c.OnHTML("div.alt.p-1 div.row.no-gutters div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "col-6 pbreak" {
			prevTitle = e.ChildText("a")
			perveMangadexID = e.ChildAttr("a", "href")
			mangaupdatesIDRaw := strings.Split(perveMangadexID, "series.html?id=")
			if len(mangaupdatesIDRaw) >= 2 {
				perveMangadexID = mangaupdatesIDRaw[1]
			}
		}

		if e.Attr("class") == "col-2 pl-1 pbreak" {
			mangahubTitle := mangahubTitleCostructor(prevTitle)
			mangaLastChapters := strings.Split(e.Text, "c.")
			var mangaLastChapterString string
			if len(mangaLastChapters) > 0 {
				mangaLastChapterString = mangaLastChapters[len(mangaLastChapters)-1]
			} else {
				mangaLastChapterString = ""
			}

			mangaLastChapterBreaked := strings.Split(mangaLastChapterString, "-")
			if len(mangaLastChapterBreaked) > 0 {
				mangaLastChapterString = mangaLastChapterBreaked[len(mangaLastChapterBreaked)-1]
			}

			mangaLastChapterString = strings.Replace(mangaLastChapterString, "a", "", -1)
			mangaLastChapterString = strings.Replace(mangaLastChapterString, "b", "", -1)
			mangaLastChapterString = strings.Replace(mangaLastChapterString, "c", "", -1)
			var mangaLastChapter float64
			mangaLastChapter, err := strconv.ParseFloat(mangaLastChapterString, 64)
			if err != nil {
				mangaLastChapter = 150
			}

			mangaData := models.MangaData{
				Title:            mangahubTitle,
				CompactTitle:     prevTitle,
				MangaLastChapter: int(mangaLastChapter),
				Weight:           idx,
				MangaUpdatesID:   perveMangadexID,
			}
			idx--
			animapuMangas = append(animapuMangas, mangaData)
		}
	})

	c.SetRequestTimeout(60 * time.Second)
	var err error
	err = c.Visit("https://www.mangaupdates.com/releases.html?page=1")
	if err != nil {
		logrus.Errorf("There is some error: %v", err)
		return animapuMangas, err
	}
	err = c.Visit("https://www.mangaupdates.com/releases.html?page=2")
	if err != nil {
		logrus.Errorf("There is some error: %v", err)
		return animapuMangas, err
	}

	return animapuMangas, nil
}

func mangahubTitleCostructor(source string) string {
	result := strings.ToLower(source)
	result = strings.Replace(result, "%", "", -1)
	result = strings.Replace(result, "'", "-", -1)
	result = strings.Replace(result, "?", "", -1)
	result = strings.Replace(result, ".", "", -1)
	result = strings.Replace(result, "&", "", -1)
	result = strings.Replace(result, ":", "", -1)
	result = strings.Replace(result, ",", "", -1)
	result = strings.Replace(result, "(", "", -1)
	result = strings.Replace(result, ")", "", -1)
	result = strings.Replace(result, "-", "", -1)
	result = strings.Replace(result, "\"", "", -1)
	result = strings.Replace(result, "  ", "-", -1)
	result = strings.Replace(result, " ", "-", -1)
	return result
}
