package scrapper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/pkg/app_cache"
)

func FetchAllAnime(c *gin.Context) map[string]string {
	animesMap := make(map[string]string)

	appCache := pkgAppCache.GetAppCache()

	res, found := appCache.Get("animepahe_map")
	if found {
		fmt.Println("FETCH FROM APP CACHE")
		return res.(map[string]string)
	}

	gColly := colly.NewCollector()

	gColly.OnHTML("div.col-12", func(e *colly.HTMLElement) {
		animeLink := e.ChildAttr("a", "href")
		animeTitle := e.ChildAttr("a", "title")
		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
		sanitizedAnimeTitle := reg.ReplaceAllString(animeTitle, "")
		animeLink = strings.ReplaceAll(animeLink, "/anime/", "")
		animesMap[sanitizedAnimeTitle] = animeLink
		fmt.Println("ANIME", animeLink)
	})

	gColly.Request("GET", "https://animepahe.com/anime", nil, nil, c.Request.Header)

	// resp, err := http.Get("https://animepahe.com/anime")
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("NORMAL HTTP:", string(body), ", ERROR:", err)

	appCache.Set("animepahe_map", animesMap, 50*time.Minute)
	return animesMap
}
