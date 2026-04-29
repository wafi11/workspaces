package services

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/core/api-gateway/config"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
)

func NewAuthRoutes(c *echo.Echo, conf *config.Config) {
	conn, err := pkg.NewGrpcConnection(&conf.AuthServiceUrl, "Auth")

	if err != nil {
		log.Printf("erorr connection auth service : %s", err.Error())
		return
	}
	svc := NewAuthService(v1.NewAuthServiceClient(conn))
	handler := NewAuthHandler(svc)
	auth := c.Group("/api/auth")

	auth.POST("/register", handler.HandleRegister)
	auth.POST("/login", handler.HandleLogin)

}
