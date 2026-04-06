package server

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/services/api-gateway/auth"
	"github.com/wafi11/workspaces/services/api-gateway/template"
	"github.com/wafi11/workspaces/services/api-gateway/user"
	"github.com/wafi11/workspaces/services/api-gateway/workspace"
	logservice "github.com/wafi11/workspaces/services/log-service"
)

func NewServer(e *echo.Echo, db *sqlx.DB, redis *redis.Client, minioClient *minio.Client, conf *config.Config,esClient *config.Client) {
	auth.RegisterRoutes(e, db, redis, conf)
	user.RegisterRoutes(e, db, redis, conf)
	template.NewTemplateRouter(e, db, redis, minioClient)
	workspace.RegisterRoutes(e, db, redis, conf, minioClient, context.Background())

	e.GET("/api/v1/logs/stream", logservice.StreamLogs(esClient))
}
