package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/internal"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
	"github.com/wafi11/workspaces/core/auth-services/tests/mocks"
)

// ─── Register ─────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	req := &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	}

	mockRepo.On("Register", mock.Anything, req).
		Return(&v1.RegisterResponse{Message: "Successfully Register"}, nil)

	resp, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Successfully Register", resp.Message)
	mockRepo.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"no @ symbol", "wafiexample.com"},
		{"no domain", "wafi@"},
		{"no tld", "wafi@example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &v1.RegisterRequest{
				Username: "wafi11",
				Email:    tt.email,
				Password: "Secret123",
			}

			resp, err := svc.Register(context.Background(), req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			mockRepo.AssertNotCalled(t, "Register")
		})
	}
}

func TestRegister_InvalidPassword(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"empty", "", pkg.ErrPasswordRequired},
		{"too short", "Ab1", pkg.ErrPasswordTooShort},
		{"no uppercase", "secret123", pkg.ErrPasswordWeak},
		{"no lowercase", "SECRET123", pkg.ErrPasswordWeak},
		{"no digit", "SecretAbc", pkg.ErrPasswordWeak},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &v1.RegisterRequest{
				Username: "wafi11",
				Email:    "wafi@example.com",
				Password: tt.password,
			}

			resp, err := svc.Register(context.Background(), req)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, resp)
			mockRepo.AssertNotCalled(t, "Register")
		})
	}
}

func TestRegister_EmailAlreadyExist(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	req := &v1.RegisterRequest{
		Username: "wafi11",
		Email:    "wafi@example.com",
		Password: "Secret123",
	}

	mockRepo.On("Register", mock.Anything, req).
		Return(nil, pkg.ErrEmailAlreadyExist)

	resp, err := svc.Register(context.Background(), req)

	assert.ErrorIs(t, err, pkg.ErrEmailAlreadyExist)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

// ─── Login ─────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	req := &v1.LoginRequest{
		Email:    "wafi@example.com",
		Password: "Secret123",
	}

	mockRepo.On("Login", mock.Anything, req).
		Return(&v1.LoginResponse{
			AccessToken:  "access-token-xyz",
			RefreshToken: "refresh-token-xyz",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}, nil)

	resp, err := svc.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "Bearer", resp.TokenType)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	req := &v1.LoginRequest{
		Email:    "bukan-email",
		Password: "Secret123",
	}

	resp, err := svc.Login(context.Background(), req)

	assert.ErrorIs(t, err, pkg.ErrEmailInvalid)
	assert.Nil(t, resp)
	mockRepo.AssertNotCalled(t, "Login")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := internal.NewService(mockRepo)

	req := &v1.LoginRequest{
		Email:    "wafi@example.com",
		Password: "WrongPass1",
	}

	mockRepo.On("Login", mock.Anything, req).
		Return(nil, pkg.ErrInvalidCredentials)

	resp, err := svc.Login(context.Background(), req)

	assert.ErrorIs(t, err, pkg.ErrInvalidCredentials)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}
