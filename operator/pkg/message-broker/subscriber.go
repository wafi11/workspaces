package messagebroker

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	workspacev1 "github.com/wafi11/workspace-operator/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type Subscriber struct {
	redis    *redis.Client
	jobQueue chan<- *workspacev1.WorkspaceEnvelope
}

func NewSubscriber(redis *redis.Client, jobQueue chan<- *workspacev1.WorkspaceEnvelope) *Subscriber {
	return &Subscriber{
		redis:    redis,
		jobQueue: jobQueue,
	}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.redis.Subscribe(ctx, ChannelOperator)
	defer pubsub.Close()

	log.Println("[subscriber] listening on channel:", ChannelOperator)

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
			s.handle(ctx, msg.Payload)
		}
	}
}

func (s *Subscriber) handle(ctx context.Context, payload string) {
	var event workspacev1.WorkspaceEnvelope
	if err := proto.Unmarshal([]byte(payload), &event); err != nil {
		log.Println("[subscriber] invalid payload:", err)
		return
	}

	switch p := event.Payload.(type) {
	case *workspacev1.WorkspaceEnvelope_Create:
		log.Printf("[subscriber] create workspace_id=%s", p.Create.Identity.WorkspaceId)
	case *workspacev1.WorkspaceEnvelope_Add:
		log.Printf("[subscriber] add pod workspace_id=%s", p.Add.Identity.WorkspaceId)
	case *workspacev1.WorkspaceEnvelope_Delete:
		log.Printf("[subscriber] delete workspace_id=%s", p.Delete.Identity.WorkspaceId)
	case *workspacev1.WorkspaceEnvelope_Stop:
		log.Printf("[subscriber] stop workspace_id=%s", p.Stop.Identity.WorkspaceId)
	case *workspacev1.WorkspaceEnvelope_Start:
		log.Printf("[subscriber] start workspace_id=%s", p.Start.Identity.WorkspaceId)
	case *workspacev1.WorkspaceEnvelope_CreatePort:
		log.Printf("[subscriber] create port on workspace_name=%s", p.CreatePort.WorkspaceName)
	case *workspacev1.WorkspaceEnvelope_DeletePort:
		log.Printf("[subscriber] delete port on workspace_name=%s", p.DeletePort.WorkspaceName)
	default:
		log.Println("[subscriber] unknown payload type")
		return
	}

	// non-blocking send ke jobQueue
	select {
	case s.jobQueue <- &event:
	case <-ctx.Done():
		log.Println("[subscriber] context cancelled, dropping event")
	}
}
