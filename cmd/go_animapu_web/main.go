package main

import (
	"fmt"

	"github.com/umarkotak/go-animapu/internal/app"
)

func main() {
	fmt.Println("Welcome to go-animapu WEB")

	app.Init()
	app.Start()
}
