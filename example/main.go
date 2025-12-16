package main

import (
	"fmt"
	"log"

	"github.com/lubasinkal/gniphyl"
)

func main() {
	// Load configuration
	config, err := gniphyl.LoadExtensionsConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Get file category
	category := gniphyl.GetFileCategory(".jpg", config)
	fmt.Println("Category for .jpg:", category)

	// Get supported categories
	categories := gniphyl.GetSupportedCategories(config)
	fmt.Println("Supported categories:", categories)

	// Get extensions for a specific category

exts, err := gniphyl.GetExtensionsForCategory("images", config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image extensions:", exts)
}