package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	authservices "github.com/wafi11/workspaces/services/auth-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redisClient *redis.Client, conf *config.Config) {

	repo := authservices.NewRepository(db, redisClient, conf)
	svc := authservices.NewServices(repo, conf)
	h := NewHandler(svc)

	auth := e.Group("/api/v1/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.RefreshToken)
}
