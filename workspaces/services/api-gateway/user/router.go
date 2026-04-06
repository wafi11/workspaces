package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/middlewares"
	userservices "github.com/wafi11/workspaces/services/user-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redis *redis.Client, conf *config.Config) {

	repo := userservices.NewRepository(db, redis)
	svc := userservices.NewServices(repo)
	h := NewHandler(svc)

	user := e.Group("/api/v1/users", middlewares.AuthMiddleware(conf))
	user.GET("", h.GetProfile)
	user.PUT("/", h.UpdateUser)
	user.PUT("/password", h.ChangePassword)
	user.GET("/sessions", h.GetUserSessions)
	user.DELETE("/sessions/:session_id", h.RevokeSession)
}
