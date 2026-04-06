package config

import (
	"log"

	"github.com/minio/minio-go/v7"
)

func NewMinio() *minio.Client {
	storage, err := minio.New("192.168.1.10:9000", &minio.Options{})

	if err != nil {
		log.Printf("failed to conn minio : %s", err.Error())
		return nil
	}

	return storage
}
