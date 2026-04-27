package internal

import (
	"context"

	"github.com/wafi11/workspaces/core/auth-services/config"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

func (repo *Repository) GenerateToken(c context.Context, req pkg.GenerateTokenReq) (*pkg.GenerateTokenRes, error) {
	accessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    req.UserID,
		Role:      req.Role,
		Exp:       1,
		SessionID: req.SessionID,
		TokenName: config.CookieAccessTokenName,
	}, repo.config)
	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	refreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    req.UserID,
		Role:      req.Role,
		SessionID: req.SessionID,
		Exp:       24,
		TokenName: config.CookieRefreshTokenName,
	}, repo.config)
	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	return &pkg.GenerateTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
