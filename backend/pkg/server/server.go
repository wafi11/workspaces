package server

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/wafi11/workspaces/config"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/websocket"
	"github.com/wafi11/workspaces/services/api-gateway/auth"
	"github.com/wafi11/workspaces/services/api-gateway/template"
	"github.com/wafi11/workspaces/services/api-gateway/user"
	"github.com/wafi11/workspaces/services/api-gateway/workspace"
)

func NewServer(e *echo.Echo, db *sqlx.DB, redis *config.RedisConnection, minioClient *minio.Client, conf *config.Config, esClient *config.Client, sub *messagebroker.Subscriber, jobQueue <-chan *proto.WorkspaceEnvelope, hub *websocket.Hub,mux *asynq.ServeMux) {
	auth.RegisterRoutes(e, db, redis.Redis, conf)
	user.RegisterRoutes(e, db, redis.Redis, conf)
	template.NewTemplateRouter(e, db, redis.Redis, minioClient,conf)
	workspace.RegisterRoutes(e, db, redis, conf, minioClient, context.Background(), sub, jobQueue, hub,mux)
}
