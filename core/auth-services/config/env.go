package config

import (
	"os"

	"github.com/joho/godotenv"
)

type SSOGithubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type SSOGoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type RedisConfig struct {
	RedisUrl      string
	RedisUsername string
	RedisPassword string
}

type Config struct {
	JWT_SECRET   string
	DB_URL       string
	Port         string
	Redis        *RedisConfig
	Github       *SSOGithubConfig
	Google       *SSOGoogleConfig
	FRONTEND_URL string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		DB_URL:       os.Getenv("DB_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("PORT"),
		FRONTEND_URL: os.Getenv("FRONTEND_URL"),
		Redis: &RedisConfig{
			RedisUrl:      os.Getenv("REDIS_URL"),
			RedisUsername: os.Getenv("REDIS_USERNAME"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		},
		Github: &SSOGithubConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
		Google: &SSOGoogleConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}
}
