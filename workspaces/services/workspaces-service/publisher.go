package workspaceservice

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

const WorkspaceEventChannel = "workspace:events"

type WorkspaceJob struct {
	WorkspaceId string
	UserId      string
	TemplateId  string
	Username    string
	Namespace   string
	Image       string
	EnvVars     map[string]any
	Action      JobAction
}

type JobAction string

const (
	JobCreate JobAction = "create"
	JobDelete JobAction = "delete"
)

func PublishEvent(ctx context.Context, rdb *redis.Client, event WorkspaceJob) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Println("[publisher] marshal error:", err)
		return
	}
	if err := rdb.Publish(ctx, WorkspaceEventChannel, payload).Err(); err != nil {
		log.Println("[publisher] publish error:", err)
	}
}
