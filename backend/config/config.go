package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL       string
	REDISURL    string
	JWT_SECRET  string
	K8S_CONFIG  string
	ELASTIC_URL string
	Port        string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		DBURL:       os.Getenv("DB_URL"),
		REDISURL:    os.Getenv("REDIS_URL"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
		ELASTIC_URL: os.Getenv("ELASTIC_URL"),
		Port:        "8080",
		K8S_CONFIG:  os.Getenv("K8S_CONFIG"),
	}
}
