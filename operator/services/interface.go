package services

import (
	"context"
)

type IK8SClient interface {
	CreateNamespace(ctx context.Context, userId string) error
	DeleteNamespace(ctx context.Context, userId string) error
	CreateResourceQuota(ctx context.Context, userId string, quota QuotaConfig) error
	UpdateResourceQuota(ctx context.Context, userId string, quota QuotaConfig) error
	SetupRBAC(ctx context.Context, userId string) error
	createServiceAccount(ctx context.Context, userID string) error
	createRole(ctx context.Context, userID string) error
	createRoleBinding(ctx context.Context, userID string) error
	DeleteRBAC(ctx context.Context, userID string) error
}

type EventType string

const (
	EventCreateWorkspace EventType = "workspace.create"
	EventDeleteWorkspace EventType = "workspace.delete"
	EventStopWorkspace   EventType = "workspace.stop"
)

const WorkspaceEventChannel = "workspace:events"

func getEnvString(envVars map[string]any, key string) string {
	v, _ := envVars[key].(string)
	return v
}
