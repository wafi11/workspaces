package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/middlewares"
	authservices "github.com/wafi11/workspaces/services/auth-service"
	userservices "github.com/wafi11/workspaces/services/user-service"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, redisClient *redis.Client, conf *config.Config) {

	repo := authservices.NewRepository(db, redisClient, conf)
	userRepo := userservices.NewRepository(db, redisClient)
	svc := authservices.NewServices(repo, conf, userRepo)
	h := NewHandler(svc)

	auth := e.Group("/api/v1/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.GET("/login", h.LoginPage)
	auth.POST("/logout", h.Logout)
	auth.GET("/validate",h.Validate)
	auth.POST("/refresh", h.RefreshToken)
	
	protected := auth.Group("",middlewares.AuthMiddleware(conf))
	protected.POST("/pat", h.CreatePAT)
	protected.GET("/pat", h.GetAllPAT)
	protected.DELETE("/pat/:pat_id", h.DeletePAT)
}
