package pkg

import (
	"context"

	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
)

type IRepository interface {
	Register(c context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error)
	Login(c context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error)
	ValidateToken(c context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error)
	RefreshToken(c context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error)
}

type GenerateTokenReq struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
}

type GenerateTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
