package network

import (
	"github.com/gin-gonic/gin"

	cPing "github.com/umarkotak/go-animapu/internal/controller"
	cAnalytic "github.com/umarkotak/go-animapu/internal/controller/analytic"
	cChats "github.com/umarkotak/go-animapu/internal/controller/chats"
	cClips "github.com/umarkotak/go-animapu/internal/controller/clips"
	cManga "github.com/umarkotak/go-animapu/internal/controller/manga"
	cUser "github.com/umarkotak/go-animapu/internal/controller/user"
	"github.com/umarkotak/go-animapu/internal/models"
)

// RouterStart this is entry porint for all http request to go-animapu web
func RouterStart(port string) {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", cPing.GetPing)

	// mangas
	router.GET("/mangas", cManga.GetManga)
	router.GET("/mangas/update", cManga.UpdateManga)
	router.GET("/mangas/firebase", cManga.GetMangaFirebase)
	router.GET("/mangas/firebase/update", cManga.UpdateMangaFirebase)
	router.GET("/mangas/search_v1", cManga.GetMangaSearch)
	router.GET("/mangas/todays_v1", cManga.GetMangaTodays)

	// users
	router.POST("/users/register", cUser.RegisterUserFirebase)
	router.OPTIONS("/users/register", cUser.SkipCors)
	router.POST("/users/login", cUser.LoginUser)
	router.OPTIONS("/users/login", cUser.SkipCors)
	router.POST("/users/read_histories", cUser.LogReadHistories)
	router.OPTIONS("/users/read_histories", cUser.SkipCors)
	router.GET("/users/detail", cUser.GetDetailFirebase)
	router.OPTIONS("/users/detail", cUser.SkipCors)

	// users analytic
	router.POST("/users/analytic_v1", cAnalytic.PostUserAnalyticV1)
	router.OPTIONS("/users/analytic_v1", cAnalytic.SkipCors)

	// clip
	router.GET("/clips", cClips.GetClips)
	router.POST("/clips", cClips.CreateClip)
	router.OPTIONS("/clips", cUser.SkipCors)

	hub := models.NewHub()
	go hub.Run()

	router.GET("/chats_v1", func(c *gin.Context) {
		cChats.GetChatsConnection(c, hub)
	})

	router.Run(":" + port)
}
