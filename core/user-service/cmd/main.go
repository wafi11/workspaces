package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wafi11/workspaces/core/user-service/config"
	v1 "github.com/wafi11/workspaces/core/user-service/gen/v1"
	"github.com/wafi11/workspaces/core/user-service/internal"
	"github.com/wafi11/workspaces/core/user-service/pkg"
	"google.golang.org/grpc"
)

func main() {
	conf := config.Load()

	conn, err := config.DBConn(conf.DB_URL)

	if err != nil {
		fmt.Printf("database : %s", err.Error())
		return
	}
	storageConn, err := pkg.NewGrpcConnection(&conf.StorageServiceUrl, "Storage")
	if err != nil {
		log.Fatalf("storage-service connection gagal: %v", err)
	}
	storageSvc := v1.NewStorageServiceClient(storageConn)

	repo := internal.NewRepository(conn, storageSvc)
	svc := internal.NewService(repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", conf.Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("user-service running on port : %s", conf.Port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
