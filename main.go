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
