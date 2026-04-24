package authservices

import (
	"context"
	"fmt"
	"strconv"

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

func (Services *Services) DeletePAT(c context.Context, PatId, userId string) error {
	return Services.repo.DeletePAT(c, PatId, userId)
}

func (services *Services) CreatePAT(c context.Context, req *CreatePATRequest) (*CreatePATResponse, error) {
	return services.repo.CreatePAT(c, req)
}
func (services *Services) Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	return services.repo.Logout(c, req)
}
func (services *Services) RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return services.repo.RefreshToken(c, req)
}

func (services *Services) Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	return services.repo.Register(c, req, UserProvidersLocal)
}

func (service *Services) Validate(c context.Context, req string) (bool, error) {
	return service.repo.Validate(c, req)
}

func (service *Services) GetAllPAT(c context.Context, userID string) ([]Pat, error) {
	return service.repo.GetAllPAT(c, userID)
}

func (s *Services) LoginWithGithub(ctx context.Context, accessToken string) (*LoginResponse, error) {
	githubUser, err := fetchGithubUser(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github user: %w", err)
	}
	return s.repo.LoginOrRegisterOAuth(ctx, &OAuthRequest{
		Email:      githubUser.Email,
		Username:   githubUser.Login,
		AvatarURL:  githubUser.AvatarURL,
		Provider:   UserProvidersGithub,
		ProviderId: strconv.Itoa(githubUser.ID),
	})
}

func (s *Services) LoginWithGoogle(ctx context.Context, accessToken string) (*LoginResponse, error) {
	googleUser, err := fetchGoogleUser(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch google user: %w", err)
	}
	return s.repo.LoginOrRegisterOAuth(ctx, &OAuthRequest{
		Email:     googleUser.Email,
		Username:  googleUser.Name,
		AvatarURL: googleUser.Picture,
		Provider:  UserProvidersGoogle,
	})
}
