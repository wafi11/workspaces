package authservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/config"
)

type Repository struct {
	db    *sqlx.DB
	redis *redis.Client
	conf  *config.Config
}

func NewRepository(db *sqlx.DB, redis *redis.Client, conf *config.Config) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
		conf:  conf,
	}
}

func (repo *Repository) Register(c context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	var userId string

	query := `
		INSERT INTO users (id,username, email, password)
		VALUES ($1, $2, $3,$4)
		RETURNING id
	`

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	err = repo.db.QueryRowContext(c, query, uuid.New(), req.Username, req.Email, hashedPassword).Scan(&userId)
	if err != nil {
		// wrap error asli untuk logging, tapi return pesan generic ke client
		return nil, fmt.Errorf("username or email already registered: %w", err)
	}

	return &RegisterResponse{
		UserId:  userId,
		Message: "Successfully created user",
	}, nil
}

func (repo *Repository) Login(c context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	var id, password, username string

	query := `
		SELECT id, username, password FROM users WHERE email = $1
	`

	err := repo.db.QueryRowContext(c, query, req.Email).Scan(&id, &username, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !VerifyPassword(password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// generate access token
	accessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    id,
		Username:  username,
		Exp:       1,
		TokenName: "access_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// generate refresh token
	refreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    id,
		Username:  username,
		Exp:       24, // 7 hari
		TokenName: "refresh_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	redisKey := fmt.Sprintf("refresh_token:%s", id)
	err = repo.redis.Set(c, redisKey, refreshToken, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// simpan session ke database
	sessionId := uuid.New().String()
	sessionQuery := `
		INSERT INTO sessions (id, user_id, is_active, user_agent, ip_address)
		VALUES ($1, $2, true, $3, $4)
	`
	_, err = repo.db.ExecContext(c, sessionQuery, sessionId, id, userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       id,
		SessionId:    sessionId,
	}, nil
}

func (repo *Repository) RefreshToken(c context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// validasi refresh token
	claims, err := config.ValidationToken(req.RefreshToken, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	// cek apakah refresh token masih ada di redis
	redisKey := fmt.Sprintf("refresh_token:%s", claims.UserID)
	storedToken, err := repo.redis.Get(c, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}

	if storedToken != req.RefreshToken {
		return nil, fmt.Errorf("refresh token mismatch")
	}

	// generate access token baru
	newAccessToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Exp:       1,
		TokenName: "access_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// rotate refresh token (best practice: tiap refresh, token baru diterbitkan)
	newRefreshToken, err := config.GenerateToken(c, &config.TokenRequest{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Exp:       24,
		TokenName: "refresh_token",
	}, repo.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// update redis dengan refresh token baru
	err = repo.redis.Set(c, redisKey, newRefreshToken, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (repo *Repository) Logout(c context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// nonaktifkan session di database
	query := `
		UPDATE sessions SET is_active = false WHERE id = $1
	`
	result, err := repo.db.ExecContext(c, query, req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate session: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("session not found")
	}

	// ambil user_id dari session untuk hapus redis
	var userId string
	err = repo.db.QueryRowContext(c,
		`SELECT user_id FROM sessions WHERE id = $1`, req.SessionId,
	).Scan(&userId)

	// hapus refresh token dari redis
	if err == nil {
		redisKey := fmt.Sprintf("refresh_token:%s", userId)
		repo.redis.Del(c, redisKey)
	}

	return &LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}
