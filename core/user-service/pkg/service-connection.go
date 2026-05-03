package pkg

import (
	"log"

	"google.golang.org/grpc"
)

func NewGrpcConnection(addr *string, service string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully Connect to %s Service", service)

	return conn, err
}
