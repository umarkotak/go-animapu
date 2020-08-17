package main

import (
	"fmt"
	"log"
	"os"

	"github.com/umarkotak/go-animapu/internal/pkg/network"
)

func main() {
	fmt.Println("Welcome to go-animapu WEB")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	network.RouterStart(port)
}
