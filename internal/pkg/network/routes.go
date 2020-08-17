package network

import (
	"github.com/gin-gonic/gin"

	cPing "github.com/umarkotak/go-animapu/internal/controller"
	cManga "github.com/umarkotak/go-animapu/internal/controller/manga"
)

// RouterStart this is entry porint for all http request to go-animapu web
func RouterStart(port string) {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", cPing.GetPing)
	router.GET("/mangas", cManga.GetManga)
	router.GET("/mangas/update", cManga.UpdateManga)

	router.Run(":" + port)
}
