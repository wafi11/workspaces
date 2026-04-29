package internal

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	v1 "github.com/wafi11/workspaces/core/user-service/gen/v1"
	"github.com/wafi11/workspaces/core/user-service/pkg"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) *Repository {
	return &Repository{
		DB: DB,
	}
}

func (repo *Repository) GetProfile(c context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	var user v1.User

	querySession := `
		select user_id from sessions where id = $1
	`

	err := repo.DB.DB.QueryRowContext(c, querySession, req.SessionId).Scan(&user.Id)

	if err != nil {
		log.Printf("error : %s", err.Error())
		return nil, pkg.ErrSessionExpired
	}

	queryUser := `
		select name,username,email,role,avatar_url from users where id = $1
	`

	err = repo.DB.DB.QueryRowContext(c, queryUser, user.Id).Scan(&user.Name, &user.Username, &user.Email, &user.Role, &user.AvatarUrl)

	if err != nil {
		log.Printf("error : %s", err.Error())

		return nil, pkg.ErrUserNotFound
	}

	return &v1.GetProfileResponse{
		User: &user,
	}, nil
}

// func (repo *Repository)
