package config

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)


type RedisConnection struct {
	Asynq *asynq.Client
	Redis *redis.Client
	Inceptor *asynq.Inspector
}

func RedisConnecion(ctx context.Context, redisUrl, redisPassword string) *RedisConnection {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})

	redisConnOpt := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisUrl,
		Password: redisPassword,
		DB: 0,
	})

	redisIncept := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: redisUrl,
		Password: redisPassword,
		DB: 0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("failed to connect redis : %s", err.Error())
		return nil
	}

	return &RedisConnection{
		Asynq: redisConnOpt,
		Redis: rdb,
		Inceptor: redisIncept,
	}
}