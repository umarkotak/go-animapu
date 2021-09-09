package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/lib/app_cache"
	"github.com/umarkotak/go-animapu/internal/lib/network"
)

func main() {
	fmt.Println("Welcome to go-animapu WEB")

	initBaseConfiguration()
	pkgAppCache.InitAppCache()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	network.RouterStart(port)
}

func initBaseConfiguration() {
	godotenv.Load(".env")
}
