package authservices

import (
	"context"
	"fmt"

	"github.com/wafi11/workspaces/config"
)

type GenerateTokenRes struct {
	AccessToken string  `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


type GenerateTokenReq struct {
	UserID string `json:"user_id"`
	Role string  `json:"role"`
	SessionID  string  `json:"session_id"`
}

func (repo *Repository) GenerateToken(c context.Context,req GenerateTokenReq)  (*GenerateTokenRes,error) {
	accessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    req.UserID,
		Role:      req.Role,
		Exp:       1,
		SessionID: req.SessionID,
		TokenName: config.CookieAccessTokenName,
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    req.UserID,
		Role:      req.Role,
		SessionID: req.SessionID,
		Exp:       24,
		TokenName: config.CookieRefreshTokenName,
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &GenerateTokenRes{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},nil
}