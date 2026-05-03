package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wafi11/workspaces/core/storage-service/config"
	v1 "github.com/wafi11/workspaces/core/storage-service/gen/v1"
	"github.com/wafi11/workspaces/core/storage-service/internal"
	"google.golang.org/grpc"
)

func main() {
	conf := config.Load()
	minio, _ := config.MinioInit(conf)

	svc := internal.NewService(minio)
	lis, err := net.Listen("tcp", conf.StorageUrl)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("storage-service running on : %s", conf.StorageUrl)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterStorageServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
