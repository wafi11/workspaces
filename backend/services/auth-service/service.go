package authservices

import (
	"context"

	"github.com/wafi11/workspaces/config"
	userservices "github.com/wafi11/workspaces/services/user-service"
)

type Services struct {
	repo     *Repository
	conf     *config.Config
	userRepo userservices.UserRepository
}

func NewServices(repo *Repository, conf *config.Config, userRepo userservices.UserRepository) *Services {
	return &Services{
		repo:     repo,
		conf:     conf,
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
	return services.repo.Register(c, req)
}

func (service *Services)  Validate(c context.Context,req string) (bool,error){
	return service.repo.Validate(c,req)
}
