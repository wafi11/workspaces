package services

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/core/api-gateway/config"
	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
	"github.com/wafi11/workspaces/core/api-gateway/pkg"
)

func NewUserRoutes(c *echo.Echo, conf *config.Config) {
	authConn, err := pkg.NewGrpcConnection(&conf.AuthServiceUrl, "Auth")

	if err != nil {
		log.Printf("erorr connection auth service : %s", err.Error())
		return
	}
	userConn, err := pkg.NewGrpcConnection(&conf.UserServiceUrl, "User")
	if err != nil {
		log.Printf("erorr connection user service : %s", err.Error())
		return
	}
	svc := NewAuthService(v1.NewAuthServiceClient(authConn))
	userSvc := NewUserService(v1.NewUserServiceClient(userConn))
	handler := NewUserHandler(userSvc)

	user := c.Group("/api/user", svc.AuthMiddleware(conf))

	user.GET("/profile", handler.GetProfile)

}
