package config

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioInit(conf *Config) (*minio.Client, error) {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(conf.MinioUsername, conf.MinioPassword, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect minio : %s", err.Error())
	}

	return minioClient, nil
}
