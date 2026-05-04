package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(conf *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisUrl,
		Password: conf.RedisPassword,
		DB:       0,
	})

	if rdb == nil {
		fmt.Printf("failed to connect redis")
		return nil, fmt.Errorf("failed connect redis")
	}

	fmt.Printf("Successfully Connect Redis")

	return rdb, nil

}
