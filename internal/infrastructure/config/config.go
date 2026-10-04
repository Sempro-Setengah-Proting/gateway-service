package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

type Config struct {
	Port string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return &Config{
		Port: port,
	}
}
