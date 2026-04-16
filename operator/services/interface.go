package services

import (
	"context"

	messagebroker "github.com/wafi11/workspace-operator/pkg/message-broker"
)

type IK8SClient interface {
	CreateNamespace(ctx context.Context, userId string) error
	DeleteNamespace(ctx context.Context, userId string) error
	CreateResourceQuota(ctx context.Context, userId string, quota messagebroker.QuotaConfig) error
	UpdateResourceQuota(ctx context.Context, userId string, quota messagebroker.QuotaConfig) error
	SetupRBAC(ctx context.Context, userId string) error
	createServiceAccount(ctx context.Context, userID string) error
	createRole(ctx context.Context, userID string) error
	StopScalling(c context.Context,namespace string,pods string) error
	StartScalling(c context.Context,namespace string,pods string) error
	createRoleBinding(ctx context.Context, userID string) error
	DeleteRBAC(ctx context.Context, userID string) error
	CreatePort(ctx context.Context, userId,workspaceName string, port int) error
	ExposeToIngress(ctx context.Context, userId,workspace_name,serviceName, path string,port int32) error
	DeletePort(ctx context.Context, userId, workspaceName string, port int) error
	RemoveFromIngress(ctx context.Context, userId, workspaceName, domain string) error
}

// func getEnvString(envVars map[string]any, key string) string {
// 	v, _ := envVars[key].(string)
// 	return v
// }
