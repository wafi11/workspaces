package services

import (
	"context"
)

type IK8SClient interface {
	CreateNamespace(ctx context.Context, namespace, workspaceId, userId string) error
	DeleteNamespace(ctx context.Context, namespace string) error
	CreateResourceQuota(ctx context.Context, namespace string, quota QuotaConfig) error
	UpdateResourceQuota(ctx context.Context, namespace string, quota QuotaConfig) error
	SetupRBAC(ctx context.Context, namespace, userId string) error
	createServiceAccount(ctx context.Context, namespace, userID string) error
	createRole(ctx context.Context, namespace string) error
	createRoleBinding(ctx context.Context, namespace, userID string) error
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
