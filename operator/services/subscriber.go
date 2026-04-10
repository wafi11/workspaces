package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	redis    *redis.Client
	jobQueue chan WorkspaceJob
}

func NewSubscriber(redis *redis.Client, jobQueue chan WorkspaceJob) *Subscriber {
	return &Subscriber{
		redis:    redis,
		jobQueue: jobQueue,
	}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.redis.Subscribe(ctx, WorkspaceEventChannel)
	defer pubsub.Close()

	log.Println("[subscriber] listening on channel:", WorkspaceEventChannel)

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			log.Println("[subscriber] shutting down")
			return

		case msg, ok := <-ch:
			if !ok {
				log.Println("[subscriber] channel closed")
				return
			}
			s.handle(msg.Payload)
		}
	}
}

func (s *Subscriber) handle(payload string) {
	var event WorkspaceJob
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Println("[subscriber] invalid payload:", err)
		return
	}

	log.Printf("[subscriber] received event: %s workspace=%s", event.Action, event.WorkspaceId)
	s.jobQueue <- event
}
