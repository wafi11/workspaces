package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wafi11/workspaces/core/api-gateway/config"
	"github.com/wafi11/workspaces/core/api-gateway/services"
)

func main() {
	conf := config.Load()

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{conf.FRONTEND_URL},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS", "PUT"},
		AllowCredentials: true,
		MaxAge:           2000,
	}))

	NewServer(e, *conf)
	log.Printf("starting backend workspace on port %s", conf.Port)
	if err := e.Start(":" + conf.Port); err != nil {
		log.Println("server stopped:", err)
	}

}

func NewServer(c *echo.Echo, conf config.Config) {

	services.NewAuthRoutes(c, &conf)
}
