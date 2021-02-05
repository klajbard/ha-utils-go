package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Fatal("Error loading .ENV file")
	}
}
