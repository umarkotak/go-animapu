main.go

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	RunManga()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", getRoot)
	router.GET("/ping", getPing)
	router.GET("/mangas", getMangas)
	router.GET("/update_mangas", getUpdateMangas)

	router.Run(":" + port)
}

func getRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"info":  "animap backend go",
		"owner": "animap",
	})
}

func getPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getMangas(c *gin.Context) {
	jsonFile, err := os.Open("data/mangas.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	c.JSON(200, result)
}

// type MangaDB struct {
// 	MangaDB map[string]Manga `json:"manga_db"`
// }

// type Manga struct {
// 	manga_last_chapter int    `json:"manga_last_chapter"`
// 	average_page       int    `json:"average_page"`
// 	status             string `json:"status"`
// 	image_url          string `json:"image_url"`
// 	new_added          int    `json:"new_added"`
// }

func getUpdateMangas(c *gin.Context) {
	jsonFile, err := os.Open("data/mangas_play.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	var mangaDb MangaDB
	json.Unmarshal(byteValue, &mangaDb)

	// log.Println(result)
	log.Println(result["manga_db"])
	log.Println(mangaDb)

	response := gin.H{"status": "done"}
	c.JSON(200, response)
}

```

manga.go

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
)

// MangaDB placeholder
type MangaDB struct {
	RawMangas map[string]interface{} `json:"manga_db"`
	Mangas    map[string]Manga       `json:"-"`
}

// MangaDB2 placeholder
type MangaDB2 struct {
	MangaDB map[string]*Manga `json:"manga_db"`
}

// Manga placeholder
type Manga struct {
	MangaLastChapter int    `json:"manga_last_chapter"`
	AveragePage      int    `json:"average_page"`
	Status           string `json:"status"`
	ImageURL         string `json:"image_url"`
	NewAdded         int    `json:"new_added"`
}

var mangaDbFile = "data/mangas_play.json"

// RunManga placeholder
func RunManga() {
	jsonFile, err := os.Open(mangaDbFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var rawDataMap map[string]interface{}
	json.Unmarshal([]byte(byteValue), &rawDataMap)

	var mangaDB MangaDB
	json.Unmarshal([]byte(byteValue), &mangaDB)
	mangaDB.Mangas = make(map[string]Manga)

	for k := range mangaDB.RawMangas {
		decodedManga := &Manga{}
		mapstructure.Decode(mangaDB.RawMangas[k], &decodedManga)

		jsonbody, err := json.Marshal(mangaDB.RawMangas[k])
		if err != nil {
			fmt.Println(err)
			return
		}
		manga := Manga{}
		if err := json.Unmarshal(jsonbody, &manga); err != nil {
			fmt.Println(err)
			return
		}
		mangaDB.Mangas[k] = manga
	}

	// var keys []string
	// for k := range mangaDB.Mangas {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)
	// for _, v := range keys {
	// 	fmt.Println(v, mangaDB.Mangas[v])
	// }

	// jsonManga, e := json.Marshal(mangaDB)
	// if e != nil {
	// 	fmt.Println("error", err)
	// }
	// os.Stdout.Write(jsonManga)

	var mangaDB2 MangaDB2
	json.Unmarshal([]byte(byteValue), &mangaDB2)
	mangaDB2.MangaDB["-- select manga title --"].AveragePage = 50
	fmt.Println(mangaDB2.MangaDB["-- select manga title --"].AveragePage)
	jsonManga2, _ := json.Marshal(mangaDB2)
	os.Stdout.Write(jsonManga2)

	ioutil.WriteFile("data/manga_play2.json", jsonManga2, 0644)
}
```