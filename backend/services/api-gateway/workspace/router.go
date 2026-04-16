package workspace

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/wafi11/workspaces/config"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/middlewares"
	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/websocket"
	workspaceservice "github.com/wafi11/workspaces/services/workspaces-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redis *config.RedisConnection, conf *config.Config, minioClient *minio.Client, ctx context.Context, sub *messagebroker.Subscriber, jobQueue <-chan *proto.WorkspaceEnvelope, hub *websocket.Hub,mux *asynq.ServeMux) {

	repo := workspaceservice.NewRepository(db, redis,hub)
	svc := workspaceservice.NewService(repo, jobQueue, hub)
	h := NewHandler(svc)
	mux.HandleFunc(string(messagebroker.EventStopWorkspace), repo.HandleStopWorkspace())

	

	ws := e.Group("/api/v1/workspaces", middlewares.AuthMiddleware(conf))
	ws.POST("", h.Create)
	ws.GET("", h.ListWorkspaces)
	ws.GET("/user", h.ListWorkspaceByUserId)
	ws.GET("/form", h.ListWorkspaceForm)
	ws.GET("/:id", h.DetailsWorkspaces)
	ws.PATCH("/:workspace_id/start", h.StartWorkspaces)
	ws.PATCH("/:workspace_id/stop", h.StopWorkspaces)
	ws.PATCH("/:workspace_id/paused", h.PausedWorkspaces)
	ws.PATCH("/:workspace_id/resumed", h.ResumedWorkspaces)
	ws.POST("/:workspace_id/port", h.CreateWorkspacePort)
	ws.GET("/:workspace_id/port", h.ListWorkspacePorts)

	svc.StartEventConsumer(ctx)

}
