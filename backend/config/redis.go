package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func RedisConnecion(ctx context.Context, redisUrl, redisPassword string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("failed to connect redis : %s", err.Error())
		return nil
	}

	return rdb
}