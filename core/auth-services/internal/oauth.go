package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

func (r *Repository) OAuthLogin(c context.Context, req *v1.OAuthCallbackRequest) (*v1.LoginResponse, error) {
	var (
		email      string
		providerID string
		provider   string
	)

	switch req.Provider {
	case "github":
		githubToken, err := pkg.ExchangeGithubToken(c, *r.config.Github, req.Code)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		githubUser, err := pkg.FetchGithubUser(c, githubToken)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		email = githubUser.Email
		providerID = fmt.Sprintf("%d", githubUser.ID)
		provider = "github"

	case "google":
		googleToken, err := pkg.ExchangeGoogleToken(c, *r.config.Google, req.Code)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		googleUser, err := pkg.FetchGoogleUser(c, googleToken)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		email = googleUser.Email
		providerID = googleUser.ID
		provider = "google"

	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	user, err := r.GetProviderUser(c, email, provider, providerID)
	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	sessionId := uuid.New().String()
	token, err := r.GenerateToken(c, pkg.GenerateTokenReq{
		UserID:    user.UserID,
		SessionID: sessionId,
		Role:      user.Role,
	})

	err = insertSession(c, sessionId, token.RefreshToken, user.Role, r.db)

	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	return &v1.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    1,
	}, nil
}

func (r *Repository) ConnectOAuth(c context.Context, req *v1.ConnectOAuthRequest) (*v1.ConnectOAuthResponse, error) {
	var (
		providerID string
	)

	switch req.Provider {
	case "github":
		githubToken, err := pkg.ExchangeGithubToken(c, *r.config.Github, req.Code)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		githubUser, err := pkg.FetchGithubUser(c, githubToken)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		providerID = fmt.Sprintf("%d", githubUser.ID)

	case "google":
		googleToken, err := pkg.ExchangeGoogleToken(c, *r.config.Google, req.Code)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		googleUser, err := pkg.FetchGoogleUser(c, googleToken)
		if err != nil {
			return nil, pkg.ErrInvalidCredentials
		}
		providerID = googleUser.ID

	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	query := `
		INSERT INTO user_providers (user_id, provider, provider_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (provider, provider_id) DO NOTHING
	`
	_, err := r.db.DB.ExecContext(c, query, req.UserId, req.Provider, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect provider: %w", err)
	}

	return &v1.ConnectOAuthResponse{
		Message:   "provider connected successfully",
		Provider:  req.Provider,
		Connected: true,
	}, nil
}

func (r *Repository) GetProviderUser(c context.Context, email, provider, providerID string) (*pkg.ProviderResponse, error) {
	var user pkg.ProviderResponse
	query := `
		SELECT 
			u.id,
			u.role
		FROM users u
		JOIN user_providers up ON up.user_id = u.id
		WHERE u.email = $1
		  AND up.provider = $2
		  AND up.provider_id = $3
	`
	err := r.db.DB.QueryRowContext(c, query, email, provider, providerID).Scan(&user.UserID, &user.Role)
	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}
	return &user, nil
}
func (r *Repository) GetOAuthURL(c context.Context, req *v1.GetOAuthURLRequest) (*v1.GetOAuthURLResponse, error) {
	state := uuid.NewString()

	err := r.redis.SetEx(c, "oauth:state:"+state, req.Provider, 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	// build URL redirect ke GitHub
	oauthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&state=%s&scope=user:email",
		r.config.Github.ClientID,
		state,
	)

	return &v1.GetOAuthURLResponse{Url: oauthURL}, nil
}
