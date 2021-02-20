package scrapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/pkg/app_cache"
)

func FetchAllAnime() map[string]string {
	fmt.Println("INCOMING!!!")
	animesMap := make(map[string]string)

	appCache := pkgAppCache.GetAppCache()

	// res, found := appCache.Get("animepahe_map")
	// if found {
	// 	fmt.Println("FETCH FROM APP CACHE")
	// 	return res.(map[string]string)
	// }

	c := colly.NewCollector()

	c.OnHTML("div.col-12", func(e *colly.HTMLElement) {
		animeLink := e.ChildAttr("a", "href")
		animeTitle := e.ChildAttr("a", "title")
		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
		sanitizedAnimeTitle := reg.ReplaceAllString(animeTitle, "")
		animeLink = strings.ReplaceAll(animeLink, "/anime/", "")
		animesMap[sanitizedAnimeTitle] = animeLink
		fmt.Println("ANIME", animeLink)
	})

	c.Visit("https://animepahe.com/anime")

	resp, err := http.Get("https://animepahe.com/anime")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("NORMAL HTTP:", string(body), ", ERROR:", err)

	appCache.Set("animepahe_map", animesMap, 50*time.Minute)
	return animesMap
}
