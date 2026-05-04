package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Register(c context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.RegisterResponse), args.Error(1)
}

func (m *MockRepository) Login(c context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.LoginResponse), args.Error(1)
}

// tests/mocks/repository_mock.go

func (m *MockRepository) RefreshToken(c context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.RefreshTokenResponse), args.Error(1)
}

func (m *MockRepository) ValidateToken(c context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.ValidateTokenResponse), args.Error(1)
}

func (m *MockRepository) ConnectOAuth(c context.Context, req *v1.ConnectOAuthRequest) (*v1.ConnectOAuthResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.ConnectOAuthResponse), args.Error(1)
}

func (m *MockRepository) OAuthLogin(c context.Context, req *v1.OAuthCallbackRequest) (*v1.LoginResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.LoginResponse), args.Error(1)
}

func (m *MockRepository) GetProviderUser(c context.Context, email, provider, providerID string) (*pkg.ProviderResponse, error) {
	args := m.Called(c, email, provider, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkg.ProviderResponse), args.Error(1)
}
func (m *MockRepository) GetOAuthURL(c context.Context, req *v1.GetOAuthURLRequest) (*v1.GetOAuthURLResponse, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.GetOAuthURLResponse), args.Error(1)
}
