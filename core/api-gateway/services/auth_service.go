package services

import (
	"context"
	"time"

	v1 "github.com/wafi11/workspaces/core/api-gateway/gen/v1"
)

type AuthService struct {
	client v1.AuthServiceClient
}

func NewAuthService(client v1.AuthServiceClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.Register(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.Login(ctx, req)
}
