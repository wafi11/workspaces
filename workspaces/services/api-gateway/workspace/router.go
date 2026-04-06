package workspace

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/k8s"
	"github.com/wafi11/workspaces/pkg/middlewares"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redis *redis.Client, conf *config.Config, minioClient *minio.Client, k8sClient *k8s.K8sClient, ctx context.Context) {

	repo := workspaceservice.NewRepository(db, redis)
	svc := workspaceservice.NewService(repo)
	h := NewHandler(svc)

	// workerqueue.StartOperator(ctx, jobQueue, k8sClient, repo, templateRepo)

	user := e.Group("/api/v1/workspaces", middlewares.AuthMiddleware(conf))
	user.POST("", h.Create)
	user.GET("", h.ListWorkspaces)
	user.GET("/:id", h.DetailsWorkspaces)
}
