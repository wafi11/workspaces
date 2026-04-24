package authservices

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (repo *Repository) LoginOrRegisterOAuth(c context.Context, req *OAuthRequest) (*LoginResponse, error) {
	var userId, username, providerId string

	query := `SELECT id, username FROM users WHERE email = $1`
	err := repo.db.DB.QueryRowContext(c, query, req.Email).Scan(&userId, &username)
	if err != nil {
		log.Printf("[LoginOrRegisterOAuth] user not found: %s", err.Error())
		resp, err := repo.Register(c, &RegisterRequest{
			Username:   req.Username,
			Email:      req.Email,
			AvatarURL:  req.AvatarURL,
			Password:   "",
			ProviderId: req.ProviderId,
		}, req.Provider)

		if err != nil {
			return nil, fmt.Errorf("failed to register oauth user : %w", err)
		}

		return &LoginResponse{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			UserId:       resp.UserId,
		}, nil
	}

	queryProviders := `SELECT provider_id FROM providers WHERE user_id = $1 AND name = $2`
	err = repo.db.DB.QueryRowContext(c, queryProviders, userId, req.Provider).Scan(&providerId)
	if err != nil {
		log.Printf("[LoginOrRegisterOAuth] linking provider %s to existing user %s", req.Provider, userId)

		return nil, fmt.Errorf("failed to link provider: %w", err)

	}

	sessionId := uuid.New().String()
	token, err := repo.GenerateToken(c, GenerateTokenReq{
		UserID:    userId,
		Role:      "user",
		SessionID: sessionId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	sessionQuery := `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address, refresh_token)
		VALUES ($1, $2, true, $3, $4, $5)
	`
	_, err = repo.db.DB.ExecContext(c, sessionQuery,
		sessionId, userId, "", "", token.RefreshToken,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		Role:         "user",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserId:       userId,
		SessionId:    sessionId,
	}, nil
}
