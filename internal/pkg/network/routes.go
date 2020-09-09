package network

import (
	"github.com/gin-gonic/gin"

	cPing "github.com/umarkotak/go-animapu/internal/controller"
	cManga "github.com/umarkotak/go-animapu/internal/controller/manga"
	cUser "github.com/umarkotak/go-animapu/internal/controller/user"
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

	// users
	router.POST("/users/register", cUser.RegisterUserFirebase)
	router.OPTIONS("/users/register", cUser.SkipCors)
	router.POST("/users/login", cUser.LoginUser)
	router.OPTIONS("/users/login", cUser.SkipCors)
	router.POST("/users/read_histories", cUser.LogReadHistories)
	router.OPTIONS("/users/read_histories", cUser.SkipCors)
	router.GET("/users/detail", cUser.GetDetailFirebase)

	router.Run(":" + port)
}
