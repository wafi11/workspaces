package services

import (
	"context"
	"log"

	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
)

type UserService struct {
	client v1.UserServiceClient
}

func NewUserService(client v1.UserServiceClient) *UserService {
	return &UserService{
		client: client,
	}
}

func (s *UserService) GetProfile(c context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	log.Printf("request incoming")
	return s.client.GetProfile(c, req)
}
