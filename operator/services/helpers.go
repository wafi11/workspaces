package services

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	messagebroker "github.com/wafi11/workspace-operator/pkg/message-broker"
	workspacev1 "github.com/wafi11/workspace-operator/pkg/proto"
	"google.golang.org/protobuf/proto"
)

func toQuotaConfig(r *workspacev1.ResourceSpec) messagebroker.QuotaConfig {
	ensure := func(val, fallback string) string {
		if val == "" {
			return fallback
		}
		return val
	}
	return messagebroker.QuotaConfig{
		CPURequest:    ensure(r.CpuRequest, "1"),
		CPULimit:      ensure(r.CpuLimit, "1"),
		MemoryRequest: ensure(r.MemoryRequest, "1024Mi"),
		MemoryLimit:   ensure(r.MemoryLimit, "1024Mi"),
		StorageLimit:  ensure(r.StorageLimit, "5Gi"),
	}
}

func toDeployParams(id *workspacev1.WorkspaceIdentity, res *workspacev1.ResourceSpec, env map[string]string) DeployParams {
	dbName := env["DB_NAME"]
	return DeployParams{
		WS_TOKEN:         env["WS_TOKEN"],
		WS_REFRESH_TOKEN: env["WS_REFRESH_TOKEN"],
		WS_API_URL:       env["WS_API_URL"],
		Namespace:        generateNamespace(id.UserId),
		DB_NAME:          &dbName,
		User:             &id.UserId,
		Password:         id.Password,
		Name:             id.UserId,
		StorageClass:     "nfs",
		StorageSize:      res.StorageRequest,
		Replicas:         1,
		RunAsUser:        1000,
		RunAsGroup:       1000,
		FsGroup:          1000,
		CPURequest:       res.CpuRequest,
		CPULimit:         res.CpuLimit,
		MemRequest:       res.MemoryRequest,
		WsID:             "",
		MemLimit:         res.MemoryLimit,
		Username:         id.Username,
		Domain:           "wfdnstore.online",
	}
}

// publish UpdateStatus FAILED → repository subscribe → rollback DB
func publishFailed(ctx context.Context, rdb *redis.Client, envelope *workspacev1.WorkspaceEnvelope, err error) {
	workspaceId := extractWorkspaceId(envelope)

	failed := &workspacev1.WorkspaceEnvelope{
		Payload: &workspacev1.WorkspaceEnvelope_Update{
			Update: &workspacev1.UpdateStatusEvent{
				WorkspaceId: workspaceId,
				Status:      workspacev1.WorkspaceStatus_WORKSPACE_STATUS_FAILED,
				Reason:      err.Error(),
			},
		},
	}

	payload, merr := proto.Marshal(failed)
	if merr != nil {
		log.Printf("[operator] marshal failed event: %v", merr)
		return
	}

	if perr := rdb.Publish(ctx, messagebroker.WorkspaceEventChannel, payload).Err(); perr != nil {
		log.Printf("[operator] publish failed event: %v", perr)
	}
}

func extractWorkspaceId(envelope *workspacev1.WorkspaceEnvelope) string {
	switch p := envelope.Payload.(type) {
	case *workspacev1.WorkspaceEnvelope_Create:
		return p.Create.Identity.WorkspaceId
	case *workspacev1.WorkspaceEnvelope_Add:
		return p.Add.Identity.WorkspaceId
	case *workspacev1.WorkspaceEnvelope_Delete:
		return p.Delete.Identity.WorkspaceId
	case *workspacev1.WorkspaceEnvelope_Stop:
		return p.Stop.Identity.WorkspaceId
	case *workspacev1.WorkspaceEnvelope_Start:
		return p.Start.Identity.WorkspaceId
	default:
		return "unknown"
	}
}
