package main

import (
	"fmt"
	"log"
	"practice7/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env from project root
	fmt.Println("Starting!")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
