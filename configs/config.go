package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TOKEN string
}

type DbConfig struct {
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
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

func NewDBConfig(debug bool) *DbConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := DbConfig{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	}

	if debug {
		db.DB_HOST = "localhost"
		db.DB_USER = "root"
		db.DB_PASSWORD = "wjdtlr21"
	}

	return &db
}
