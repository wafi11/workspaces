package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
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
