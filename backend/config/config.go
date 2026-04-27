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
type Config struct {
	DBURL        string
	REDISURL     string
	JWT_SECRET   string
	K8S_CONFIG   string
	ELASTIC_URL  string
	Port         string
	Github       *SSOGithubConfig
	Google       *SSOGoogleConfig
	FRONTEND_URL string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		DBURL:        os.Getenv("DB_URL"),
		REDISURL:     os.Getenv("REDIS_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		ELASTIC_URL:  os.Getenv("ELASTIC_URL"),
		Port:         "8080",
		K8S_CONFIG:   os.Getenv("K8S_CONFIG"),
		FRONTEND_URL: os.Getenv("FRONTEND_URL"),
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
