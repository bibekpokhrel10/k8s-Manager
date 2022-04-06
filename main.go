package main

import (
	"log"

	"github.com/joho/godotenv"
	"wordpress.com/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router.RouterHandle()
}
