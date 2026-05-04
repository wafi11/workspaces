package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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

func (s *AuthService) GetOAuthURL(ctx context.Context, req *v1.GetOAuthURLRequest) (*v1.GetOAuthURLResponse, error) {
	state := uuid.NewString() // random state

	return s.client.GetOAuthURL(ctx, &v1.GetOAuthURLRequest{
		Provider: req.Provider,
		State:    state,
	})
}
func (s *AuthService) ConnectOAuth(ctx context.Context, in *v1.ConnectOAuthRequest) (*v1.ConnectOAuthResponse, error) {
	return s.client.ConnectOAuth(ctx, in)
}
func (s *AuthService) OAuthLogin(ctx context.Context, in *v1.OAuthCallbackRequest) (*v1.LoginResponse, error) {
	return s.client.OAuthLogin(ctx, in)
}
