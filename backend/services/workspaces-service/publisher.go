package workspaceservice

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

const WorkspaceEventChannel = "workspace:events"

type WorkspaceJob struct {
	WorkspaceId          string
	Name                 string
	UserId               string
	TemplateId           string
	Username             string
	Namespace            string
	Image                string
	EnvVars              map[string]any
	CPURequest           string
	MemoryRequest        string
	StorageRequest       string
	MemoryTerminalLimit  string
	StorageTerminalLimit string
	CpuTerminalLimit     string
	CPULimit             string
	MemoryLimit          string
	StorageLimit         string
	Action               JobAction
	Replicas             string
}

type JobAction string

const (
	JobCreate JobAction = "create"
	JobDelete JobAction = "delete"
	JobAdd    JobAction = "add"
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
