package messagebroker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	workspacev1 "github.com/wafi11/workspaces/pkg/proto"
	"google.golang.org/protobuf/proto"
)

func PublishEvent(ctx context.Context, rdb *redis.Client, event *workspacev1.WorkspaceEnvelope) error {
	payload, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("[publisher] marshal: %w", err)
	}

	if err := rdb.Publish(ctx, ChannelOperator, payload).Err(); err != nil {
		return fmt.Errorf("[publisher] publish: %w", err)
	}

	return nil
}
