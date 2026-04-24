package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wafi11/workspaces/config"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/middlewares"
	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/server"
	"github.com/wafi11/workspaces/pkg/websocket"
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

	jobQueue := make(chan *proto.WorkspaceEnvelope, 100)
	sub := messagebroker.NewSubscriber(redisClient.Redis, jobQueue)
	go sub.Start(context.Background())

	// init echo
	hub := websocket.NewHub(conf)
	e := echo.New()

	e.GET("/ws", hub.Handler, middlewares.AuthMiddleware(conf))

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://web-platform.wfdnstore.online"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS", "PUT"},
		AllowCredentials: true,
		MaxAge:           2000,
	}))

	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr:     conf.REDISURL,
		Password: "",
		DB:       0,
	}, asynq.Config{})
	mux := asynq.NewServeMux()

	k8s, err := config.NewK8sClient(conf.K8S_CONFIG)
	if err != nil {
		log.Fatalf("failed to connect k8s: %v", err)
	}

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("asynq server error: %v", err)
		}
	}()

	server.NewServer(e, database, redisClient, minio, conf, esClient, sub, jobQueue, hub, mux, k8s)
	log.Printf("starting backend workspace on port %s", conf.Port)
	if err := e.Start(":" + conf.Port); err != nil {
		log.Println("server stopped:", err)
	}

}
