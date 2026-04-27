package pkg

import (
	"log"

	"google.golang.org/grpc"
)

func NewAuthConnection(addr *string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	log.Printf("Successfully Connect to Auth Service")

	return conn, err
}
