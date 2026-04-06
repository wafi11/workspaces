package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wafi11/workspace-operator/config"
	"github.com/wafi11/workspace-operator/services"
)

func main() {
	conf := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := config.DBConnection(conf.DBURL)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	defer database.Close()

	redisClient := config.RedisConnecion(ctx, conf.REDISURL, "")
	minio := config.NewMinio()

	k8sClient, err := services.NewK8sClient(conf.K8S_CONFIG)
	if err != nil {
		log.Fatal("failed to init k8s client: ", err)
	}

	jobQueue := make(chan services.WorkspaceJob, 100)
	repo := services.NewRepository(redisClient, jobQueue, database, minio, k8sClient.DynClient, k8sClient.Mapper)

	// start operator worker
	services.StartOperator(ctx, jobQueue, k8sClient, repo)

	// start redis subscriber
	sub := services.NewSubscriber(redisClient, jobQueue)
	go sub.Start(ctx)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	cancel() // stop subscriber + operator
}
