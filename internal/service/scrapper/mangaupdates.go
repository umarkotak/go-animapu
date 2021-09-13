package scrapper

import (
	"fmt"

	"github.com/gocolly/colly"
)

type (
	MangaUpdatesResult struct {
		Title string
		Raw   []map[string]interface{}
	}
)

func GetSeries(mangaID string) ([]MangaUpdatesResult, error) {
	url := fmt.Sprintf("https://www.mangaupdates.com/series.html?id=%v", mangaID)

	mangaUpdatesResults := []MangaUpdatesResult{
		{
			Raw: []map[string]interface{}{},
		},
	}

	c := colly.NewCollector()

	c.OnHTML("span.releasestitle.tabletitle", func(e *colly.HTMLElement) {
		tempObj := map[string]interface{}{}
		tempObj["title"] = e.Text
		mangaUpdatesResults[0].Raw = append(mangaUpdatesResults[0].Raw, tempObj)
	})

	err := c.Visit(url)

	return mangaUpdatesResults, err
}
