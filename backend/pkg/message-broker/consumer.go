package messagebroker

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	protos "github.com/wafi11/workspaces/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type Subscriber struct {
	redis    *redis.Client
	jobQueue chan<- *protos.WorkspaceEnvelope
}

func NewSubscriber(redis *redis.Client, jobQueue chan<- *protos.WorkspaceEnvelope) *Subscriber {
	return &Subscriber{
		redis:    redis,
		jobQueue: jobQueue,
	}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.redis.Subscribe(ctx, ChannelBackend)
	defer pubsub.Close()

	log.Println("[subscriber] listening on channel:", ChannelBackend)

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
	var event protos.WorkspaceEnvelope

	err := proto.Unmarshal([]byte(payload), &event)
	if err != nil {
		log.Println("[subscriber] invalid payload:", err)
		return
	}

	switch e := event.Payload.(type) {
	case *protos.WorkspaceEnvelope_Update:
		log.Printf("[subscriber] received update event: %s", e.Update.WorkspaceId)
		s.jobQueue <- &event

	case *protos.WorkspaceEnvelope_Create:
		log.Println("[subscriber] ignoring create event (operator only)")
		return

	default:
		log.Println("[subscriber] unknown event type")
	}
}
