package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	StorageUrl    string
	MinioUrl      string
	MinioUsername string
	MinioPassword string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		MinioUsername: os.Getenv("MINIO_USERNAME"),
		MinioPassword: os.Getenv("MINIO_PASSWORD"),
		StorageUrl:    os.Getenv("STORAGE_URL"),
		MinioUrl:      os.Getenv("MINIO_URL"),
	}
}
