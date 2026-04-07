package authservices

import (
	"context"

	"github.com/wafi11/workspaces/config"
	"github.com/wafi11/workspaces/pkg/models"
	userservices "github.com/wafi11/workspaces/services/user-service"
)



type Services struct {
	repo *Repository
	conf *config.Config
	userRepo  userservices.UserRepository
}

func NewServices(repo *Repository, conf *config.Config, userRepo userservices.UserRepository) *Services {
	return &Services{
		repo: repo,
		conf: conf,
		userRepo: userRepo,
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
	resp,err := services.repo.Register(c, req)

	if err != nil {
		return nil, err
	}

	err = services.userRepo.CreateUserQuota(c,&models.UserQuota{
		UserID:        resp.UserId,
		MaxWorkspaces: 2,
		MaxStorageGB:  10,
		MaxRamMB:      4096,
		MaxCpuCores:   4,
	})

	return resp, nil
}
