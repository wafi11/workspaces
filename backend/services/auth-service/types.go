package authservices

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type IServices interface {
	Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error)
	Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error)
	RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Validate(c context.Context, req string) (bool, error)
	CreatePAT(c context.Context, req *CreatePATRequest) (*CreatePATResponse, error)
	DeletePAT(c context.Context, PatId, userId string) error
	GetAllPAT(c context.Context, userID string) ([]Pat, error)
	GithubOauthConfig() *oauth2.Config
	GoogleOauthConfig() *oauth2.Config
	GenerateState() string
	LoginWithGithub(ctx context.Context, accessToken string) (*LoginResponse, error)
	LoginWithGoogle(ctx context.Context, accessToken string) (*LoginResponse, error)
}

type RegisterRequest struct {
	ProviderId string `json:"provider_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	AvatarURL  string `json:"avatar_url"`
}

type RegisterResponse struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
	SessionId    string `json:"session_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	SessionId string `json:"session_id"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ValidateTokenResponse struct {
	Valid    bool
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

const (
	UserProvidersGithub string = "github"
	UserProvidersGoogle string = "google"
	UserProvidersLocal  string = "local"
)

type GithubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type OAuthRequest struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	AvatarURL  string `json:"avatar_url"`
	Provider   string `json:"provider"`
	ProviderId string `json:"provider_id"`
}
type CreatePATRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at"`
	UserId    string     `json:"user_id"`
}

type CreatePATResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type Pat struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	LastUsedAt time.Time `json:"last_used_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
