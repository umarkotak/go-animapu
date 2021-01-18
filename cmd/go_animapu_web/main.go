package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	pkgAppCache "github.com/umarkotak/go-animapu/internal/pkg/app_cache"
	"github.com/umarkotak/go-animapu/internal/pkg/network"
)

func main() {
	fmt.Println("Welcome to go-animapu WEB")

	initBaseConfiguration()
	pkgAppCache.InitAppCache()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	network.RouterStart(port)
}

func initBaseConfiguration() {
	godotenv.Load(".env")
}
