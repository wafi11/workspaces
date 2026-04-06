package server

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/k8s"
	"github.com/wafi11/workspaces/services/api-gateway/auth"
	"github.com/wafi11/workspaces/services/api-gateway/template"
	"github.com/wafi11/workspaces/services/api-gateway/user"
	"github.com/wafi11/workspaces/services/api-gateway/workspace"
	logservice "github.com/wafi11/workspaces/services/log-service"
)

func NewServer(e *echo.Echo, db *sqlx.DB, redis *redis.Client, minioClient *minio.Client, conf *config.Config, k8s *k8s.K8sClient,esClient *config.Client) {
	auth.RegisterRoutes(e, db, redis, conf)
	user.RegisterRoutes(e, db, redis, conf)
	template.NewTemplateRouter(e, db, redis, minioClient, k8s.DynClient, k8s.Mapper)
	workspace.RegisterRoutes(e, db, redis, conf, minioClient, k8s, context.Background())

	e.GET("/api/v1/logs/stream", logservice.StreamLogs(esClient))
}
