package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		db:          db,
		redisClient: redis,
	}
}

// ─── Repository Methods ───────────────────────────────────────────────────────

func (r *Repository) GetProfile(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	// 1. check cache
	if cached, err := r.getUserCache(ctx, req.UserId); err == nil {
		return &GetUserResponse{User: cached}, nil
	}

	// 2. cache miss → hit db
	var (
		profile User
	)

	query := `
		SELECT 
			id, 
			email, username,terminal_url,created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, req.UserId).
		Scan(&profile.Id, &profile.Email, &profile.Username, &profile.TerminalUrl, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		fmt.Printf("error : %s", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 3. store to cache
	r.setUserCache(ctx, req.UserId, &profile)

	return &GetUserResponse{User: &profile}, nil
}

func (r *Repository) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	query := `
		UPDATE users 
		SET username = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, username, email, created_at, updated_at
	`

	var (
		profile   User
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, query, req.Username, req.Email, req.UserId).
		Scan(&profile.Id, &profile.Username, &profile.Email, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	profile.CreatedAt = createdAt
	profile.UpdatedAt = updatedAt
	// invalidate old cache → set new cache
	r.invalidateUserCache(ctx, req.UserId)
	r.setUserCache(ctx, req.UserId, &profile)

	return &UpdateUserResponse{
		User:    &profile,
		Message: "user updated successfully",
	}, nil
}

func (r *Repository) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	var hashedPassword string
	err := r.db.QueryRowContext(ctx,
		`SELECT password FROM users WHERE id = $1`, req.UserId,
	).Scan(&hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.OldPassword)); err != nil {
		return nil, errors.New("old password is incorrect")
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`,
		newHashed, req.UserId,
	)
	if err != nil {
		return nil, err
	}

	// invalidate cache karena updated_at berubah
	r.invalidateUserCache(ctx, req.UserId)

	return &ChangePasswordResponse{Message: "password changed successfully"}, nil
}
