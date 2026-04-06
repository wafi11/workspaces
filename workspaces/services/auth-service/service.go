package authservices

import (
	"context"

	"github.com/wafi11/workspaces/config"
)

type IServices interface {
	Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error)
	Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error)
	RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error)
}

type Services struct {
	repo *Repository
	conf *config.Config
}

func NewServices(repo *Repository, conf *config.Config) *Services {
	return &Services{
		repo: repo,
		conf: conf,
	}
}

func (services *Services) Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	return services.repo.Login(c, req, userAgent, ipAddress)
}
func (services *Services) Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	return services.repo.Logout(c, req)
}
func (services *Services) RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return services.repo.RefreshToken(c, req)
}

func (services *Services) Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	return services.repo.Register(c, req)
}
