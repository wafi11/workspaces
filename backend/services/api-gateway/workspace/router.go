package workspace

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/middlewares"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redis *redis.Client, conf *config.Config, minioClient *minio.Client, ctx context.Context) {

	repo := workspaceservice.NewRepository(db, redis)
	svc := workspaceservice.NewService(repo)
	h := NewHandler(svc)

	ws := e.Group("/api/v1/workspaces", middlewares.AuthMiddleware(conf))
	ws.POST("", h.Create)
	ws.GET("", h.ListWorkspaces)
	ws.GET("/user", h.ListWorkspaceByUserId)
	ws.GET("/form", h.ListWorkspaceForm)
	ws.GET("/:id", h.DetailsWorkspaces)
	ws.PATCH("/:workspace_id/start", h.StartWorkspaces)
	ws.PATCH("/:workspace_id/stop", h.StopWorkspaces)

	ws.GET("/:workspace_id/add-ons", h.GetAddonService)
	ws.POST("/:workspace_id/add-ons", h.CreateAddonService)
}
