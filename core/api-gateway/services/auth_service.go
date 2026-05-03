package services

import (
	"context"
	"log"
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

	reg, err := s.client.Register(ctx, req)

	if err != nil {
		log.Printf("error : %s", err.Error())
		return nil, err
	}

	return reg, nil
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.Login(ctx, req)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	return s.client.RefreshToken(ctx, req)
}
