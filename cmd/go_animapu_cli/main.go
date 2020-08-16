package main

import (
	"fmt"

	"github.com/umarkotak/go-animapu/internal"
)

func main() {
	fmt.Println("Welcome to go-animapu CLI")

	internal.GetMangaFromJSON()
}
