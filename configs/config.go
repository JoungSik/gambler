package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TOKEN string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		TOKEN: os.Getenv("TOKEN"),
	}
}
