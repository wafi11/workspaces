package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthServiceUrl string
	UserServiceUrl string
	Port           string
	FRONTEND_URL   string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		AuthServiceUrl: os.Getenv("AUTH_SERVICE_URL"),
		UserServiceUrl: os.Getenv("USER_SERVICE_URL"),
		Port:           "8080",
		FRONTEND_URL:   os.Getenv("FRONTEND_URL"),
	}
}
