package network

import (
	"github.com/gin-gonic/gin"

	cPing "github.com/umarkotak/go-animapu/internal/controller"
)

// RouterStart this is entry porint for all http request to go-animapu web
func RouterStart() {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", cPing.GetPing)
}
