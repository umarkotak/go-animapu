package scrapper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/go-animapu/internal/models"
)

type (
	MangaUpdatesResult struct {
		Title string
		Raw   map[string]interface{}
	}
)

func GetSeries(mangaID string) ([]MangaUpdatesResult, error) {
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

func GetReleases() (models.MangaDB, error) {
	url := fmt.Sprintf("https://www.mangaupdates.com/releases.html")

	mangaDatas := map[string]*models.MangaData{}
	mangaDB := models.MangaDB{MangaDatas: mangaDatas}

	c := colly.NewCollector()

	prevTitle := ""
	idx := 100000

	c.OnHTML("div.alt.p-1 div.row.no-gutters div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "col-6 pbreak" {
			prevTitle = e.ChildText("a")
		}

		if e.Attr("class") == "col-2 pl-1 pbreak" {
			mangahubTitle := strings.ToLower(prevTitle)
			mangahubTitle = strings.Replace(mangahubTitle, "%", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "'", "-", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "?", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, ".", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "&", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, ":", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, ",", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "(", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, ")", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "-", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "\"", "", -1)
			mangahubTitle = strings.Replace(mangahubTitle, "  ", "-", -1)
			mangahubTitle = strings.Replace(mangahubTitle, " ", "-", -1)

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

			fmt.Println(mangahubTitle, mangaLastChapters)

			mangaData := models.MangaData{
				Title:            mangahubTitle,
				CompactTitle:     prevTitle,
				MangaLastChapter: int(mangaLastChapter),
				Weight:           idx,
			}
			_, found := mangaDatas[mangahubTitle]
			if !found {
				mangaDatas[mangahubTitle] = &mangaData
			}
			idx--
		}

	})

	c.SetRequestTimeout(60 * time.Second)
	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("There is some error: %v", err)
		return mangaDB, err
	}

	return mangaDB, nil
}
