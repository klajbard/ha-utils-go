package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Println("Error loading .ENV file")
	}
}
