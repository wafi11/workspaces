package workspaceservice

import (
	"context"
	"encoding/json"
	"log"
	"time"

	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/proto"
	authservices "github.com/wafi11/workspaces/services/auth-service"
)


func (r *Repository) publishStop(ctx context.Context, workspaceId string, res *workspaceRow) {
	messagebroker.PublishEvent(ctx, r.redis.Redis, &proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_Stop{
			Stop: &proto.StopWorkspaceEvent{
				Identity: &proto.WorkspaceIdentity{
					WorkspaceId: workspaceId,
					UserId:      res.UserId,
					Name:        res.Name,
					Namespace:   authservices.GenerateNamespace(res.UserId),
				},
			},
		},
	})
}


func (r *Repository) publishStart(ctx context.Context, workspaceId string, res *workspaceRow) {
	messagebroker.PublishEvent(ctx, r.redis.Redis, &proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_Start{
			Start: &proto.StartWorkspaceEvent{
				Identity: &proto.WorkspaceIdentity{
					WorkspaceId: workspaceId,
					UserId:      res.UserId,
					Name:        res.Name,
					Namespace:   authservices.GenerateNamespace(res.UserId),
				},
			},
		},
	})
}

func (r *Repository) scheduleAutoStop(workspaceId string, d int,typeTimeDuration time.Duration) {
	payload, err := json.Marshal(messagebroker.TaskSchedulling{
		WorkspaceID: workspaceId,
		Duration:    d,
		Status:      string(messagebroker.EventStopWorkspace),
		TypeTimeDuration: typeTimeDuration,
	})
	if err != nil {
		log.Printf("[scheduler] failed marshal payload: %v", err)
		return
	}
	if err := messagebroker.TaskStartAndStopWorkspace(string(payload), r.redis.Asynq); err != nil {
		log.Printf("[scheduler] failed enqueue stop task: %v", err)
	}
}