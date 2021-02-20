package scrapper

import (
	"fmt"
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

	header := http.Header{
		"Authority":                 []string{"animepahe.com"},
		"Cache-Control":             []string{"max-age=0"},
		"Sec-Ch-Ua":                 []string{"\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\""},
		"Sec-Ch-Ua-Mobile":          []string{"?0"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"User-Agent":                []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36"},
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Sec-Fetch-Site":            []string{"same-origin"},
		"Sec-Fetch-Mode":            []string{"navigate"},
		"Sec-Fetch-User":            []string{"?1"},
		"Sec-Fetch-Dest":            []string{"document"},
		"Referer":                   []string{"https://animepahe.com/"},
		"Accept-Language":           []string{"en-US,en;q=0.9,id;q=0.8"},
		"Cookie":                    []string{"res=720; aud=jpn; __ddg1=Stz1rnA6MfBs9wJpImMq; _pk_id.1.b28c=8f7e0f509f89a3b5.1592446632.; __ddg2=pJWtZyCeDuUbwP4j; __ddg4=cf53380212cee790925fc8067814338a; __cfduid=d9bb1bccc5cf6bff672bb105075d9102f1612401357; SERVERID=karma; latest=48016; _pk_ses.1.b28c=1; __cf_bm=858b97d693154f175cd83226295316d4cb7eb45b-1613831639-1800-AWnxpLNa/C7eEly7NFd8OvQgurHxGdnHrSaF/O2xZ63ET29ufTiWN7BYurQa5o3Xj5YlX5B5HyLIqERYAMjWXarosEc1VhDgypl0UkQ5zp4GbfZGiuZxPtkjJgq7G0WHeg==; XSRF-TOKEN=eyJpdiI6IktjZFdXWml0a1pRY0FPM1dZKzFtNVE9PSIsInZhbHVlIjoiQmlhN0dwcjVLeTFodXBQRENBbldIemI3M1lKUHl6bHV6RjNJdnZhK0tvNjYzYnV0T2xNT3dlT3VlTnRNRWlMRyIsIm1hYyI6IjQ4ZmUzMDY2NDM5NDc0YTAxNjNjOTA4NWYyOTZlMWVhM2Y1ZGIzOTI0MjcwNDQ3YzhiMmRiNzc0YTM1OWI3MWEifQ%3D%3D; laravel_session=eyJpdiI6IlwvaXoxdzB2K1hFMHQ3WDVhXC9RMkNBUT09IiwidmFsdWUiOiJOODg0K29iR2RuZlAwZWdHYlR1Umd5eW8waHQ0Y05sNGZwNXBvNEVwRnB0VWxQOVVTM2xDYmxYTHRlM3NFUkZBIiwibWFjIjoiN2IyMDYzNWM1ZjljNGY2Y2RhZjhkODkzZTBiMzUxNzg4OTVkMmYzM2JmODBkMmUxMjUxOGY4MmIzZDRkYzM0NCJ9"},
	}
	c.Request("GET", "https://animepahe.com/anime", nil, nil, header)

	// resp, err := http.Get("https://animepahe.com/anime")
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("NORMAL HTTP:", string(body), ", ERROR:", err)

	appCache.Set("animepahe_map", animesMap, 50*time.Minute)
	return animesMap
}
