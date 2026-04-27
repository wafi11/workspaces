package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wafi11/workspaces/core/auth-services/config"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/internal"
	"google.golang.org/grpc"
)

func main() {

	conf := config.Load()

	conn, err := config.DBConn(conf.DB_URL)

	if err != nil {
		fmt.Printf("database : %s", err.Error())
		return
	}

	repo := internal.NewRepository(conf, conn)
	svc := internal.NewService(repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", conf.Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("auth-service running on port : %s", conf.Port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterAuthServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
