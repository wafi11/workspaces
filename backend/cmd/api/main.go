package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/server"
)

func main() {
	// load config
	conf := config.Load()

	// init db
	database, err := config.DBConnection(conf.DBURL)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	defer database.Close()

	// init redis
	redisClient := config.RedisConnecion(context.Background(), conf.REDISURL, "")

	minio := config.NewMinio()

	esClient, err := config.NewClient(conf.ELASTIC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}


	if err != nil {
		fmt.Printf("failed to connect k8s %v+", err)
		return
	}

	// init echo
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://*"},
		AllowHeaders: []string{"Content-type", "Authorization"},
		AllowMethods: []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS", "PUT"},
		MaxAge:       2000,
	}))

	server.NewServer(e, database, redisClient, minio, conf, esClient)

	log.Printf("starting backend workspace on port %s", conf.Port)
	if err := e.Start(":" + conf.Port); err != nil {
		log.Println("server stopped:", err)
	}

}
