package publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

const WorkspaceEventChannel = "workspace:events"

type WorkspaceJob struct {
	WorkspaceId         string         `json:"workspace_id"`
	UserId              string         `json:"user_id"`
	TemplateId          string         `json:"template_id"`
	Username            string         `json:"username"`
	Name                string         `json:"name"`
	Namespace           string         `json:"namespace"`
	Image               string         `json:"image"`
	EnvVars             map[string]any `json:"env_vars"`
	CPURequest          string         `json:"cpu_request"`
	MemoryRequest       string         `json:"memory_request"`
	StorageRequest      string         `json:"storage_request"`
	MemoryTerminalLimit string         `json:"memory_terminal_limit"`
	CpuTerminalLimit    string         `json:"cpu_terminal_limit"`
	CPULimit            string         `json:"cpu_limit"`
	MemoryLimit         string         `json:"memory_limit"`
	StorageLimit        string         `json:"storage_limit"`
	Action              JobAction      `json:"action"`
	Replicas            string         `json:"replicas"`
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
