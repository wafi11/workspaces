package internal

import (
	"context"

	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

type Service struct {
	repo pkg.IRepository
	v1.UnimplementedAuthServiceServer
}

func NewService(repo pkg.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(c context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if err := pkg.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := pkg.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := pkg.ValidateUsername(req.Username); err != nil {
		return nil, err
	}
	return s.repo.Register(c, req)
}

func (s *Service) Login(c context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if err := pkg.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := pkg.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	return s.repo.Login(c, req)
}

func (s *Service) RefreshToken(c context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {

	return s.repo.RefreshToken(c, req)
}

func (s *Service) ValidateToken(c context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	return s.repo.ValidateToken(c, req)
}

func (s *Service) OAuthLogin(c context.Context, req *v1.OAuthCallbackRequest) (*v1.LoginResponse, error) {
	return s.repo.OAuthLogin(c, req)
}

func (s *Service) ConnectOAuth(c context.Context, req *v1.ConnectOAuthRequest) (*v1.ConnectOAuthResponse, error) {
	return s.repo.ConnectOAuth(c, req)
}
