package internal

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wafi11/workspaces/core/auth-services/config"
	v1 "github.com/wafi11/workspaces/core/auth-services/gen/v1"
	"github.com/wafi11/workspaces/core/auth-services/pkg"
)

type Repository struct {
	config *config.Config
	db     *sqlx.DB
}

func NewRepository(config *config.Config, db *sqlx.DB) *Repository {
	return &Repository{
		config: config,
		db:     db,
	}
}

func (r *Repository) Register(c context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	var user_id string

	query := `
		select id from users where email = $1
	`
	err := r.db.DB.QueryRowContext(c, query, req.Email).Scan(&user_id)

	if err == nil {
		return nil, pkg.ErrEmailAlreadyExist
	}

	hashedPassword, err := pkg.HashPassword(req.Password)

	if err != nil {
		return nil, pkg.ErrPasswordRequired
	}

	query = `
		insert into users (username,email,password,role) values ($1,$2,$3,'user') returning id
	`

	err = r.db.DB.QueryRowContext(c, query, req.Username, req.Email, hashedPassword).Scan(&user_id)

	if err != nil {
		return nil, pkg.ErrInternalServerError
	}

	return &v1.RegisterResponse{
		Message: "Successfully Register",
	}, nil
}

func (repo *Repository) Login(c context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	sessionId := uuid.New().String()
	var password, user_id, role string

	query := `
		select id,password,role from users where email = $1
	`

	err := repo.db.DB.QueryRowContext(c, query, req.Email).Scan(&user_id, &password, &role)

	if err != nil {
		return nil, pkg.ErrEmailInvalid
	}

	if !pkg.VerifyPassword(password, req.Password) {
		return nil, pkg.ErrInvalidCredentials
	}

	token, err := repo.GenerateToken(c, pkg.GenerateTokenReq{
		UserID:    user_id,
		Role:      role,
		SessionID: sessionId,
	})
	if err != nil {
		return nil, err
	}

	err = insertSession(c, sessionId, token.RefreshToken, user_id, repo.db)

	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	return &v1.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "cookie",
		ExpiresIn:    1,
	}, nil
}

func (repo *Repository) ValidateToken(c context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	validate, err := config.ValidationToken(req.Token, repo.config)
	var session_id string

	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	query := `
		select id from sessions where user_id = $1 and id = $2
	`

	err = repo.db.DB.QueryRowContext(c, query, validate.UserID, validate.SessionID).Scan(&session_id)

	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	return &v1.ValidateTokenResponse{
		UserId:    validate.UserID,
		SessionId: validate.SessionID,
		Role:      validate.Role,
	}, nil
}

func (repo *Repository) RefreshToken(c context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	return nil, nil
}

func insertSession(c context.Context, sessionId, refrsh_token, userId string, db *sqlx.DB) error {
	query := `
		insert into sessions (id,user_id,refresh_token) values ($1,$2,$3)
	`

	_, err := db.DB.ExecContext(c, query, sessionId, userId, refrsh_token)
	if err != nil {
		return pkg.ErrInvalidCredentials
	}

	return nil
}
