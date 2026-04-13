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
}

// func getEnvString(envVars map[string]any, key string) string {
// 	v, _ := envVars[key].(string)
// 	return v
// }
