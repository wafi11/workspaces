package k8s

import (
	"context"
)


type IK8SClient interface {
	CreateNamespace(ctx context.Context, namespace,workspaceId, userId string) error
	DeleteNamespace(ctx context.Context, namespace string) error
	CreateResourceQuota(ctx context.Context, userId string, quota QuotaConfig) error
	UpdateResourceQuota(ctx context.Context, userId string, quota QuotaConfig) error
	SetupRBAC(ctx context.Context, namespace, userId string) error
	createServiceAccount(ctx context.Context, namespace, userID string) error
	createRole(ctx context.Context, namespace string) error
	createRoleBinding(ctx context.Context, namespace, userID string) error
	DeleteRBAC(ctx context.Context, userID string) error
}

type QuotaConfig struct {
	CPULimit      string
	MemoryLimit   string
	StorageLimit  string
	CPURequest    string
	MemoryRequest string
}
