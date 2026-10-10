package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

type Config struct {
	GatewayPort        string
	IdentityServiceURL string
	JWTSecret          string
	GatewayServiceKey  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = ":8080"
	}

	identityURL := os.Getenv("IDENTITY_SERVICE_URL")
	if identityURL == "" {
		identityURL = "http://localhost:8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	gatewayServiceKey := os.Getenv("GATEWAY_SERVICE_KEY")
	if gatewayServiceKey == "" {
		log.Fatal("GATEWAY_SERVICE_KEY not set")
	}

	return &Config{
		GatewayPort:        port,
		IdentityServiceURL: identityURL,
		JWTSecret:          jwtSecret,
		GatewayServiceKey:  gatewayServiceKey,
	}

}
