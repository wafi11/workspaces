package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	v1 "github.com/wafi11/workspaces/core/user-service/gen/v1"
	"github.com/wafi11/workspaces/core/user-service/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
	DB      *sqlx.DB
	Storage v1.StorageServiceClient
}

func NewRepository(DB *sqlx.DB, Storage v1.StorageServiceClient) *Repository {
	return &Repository{
		DB:      DB,
		Storage: Storage,
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

func (r *Repository) UpdateProfile(c context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {
	var avatarURL *string // pointer biar bisa nil untuk COALESCE

	if req.AvatarBase64 != nil && len(*req.AvatarBase64) > 0 {
		raw := *req.AvatarBase64
		if idx := strings.Index(raw, ","); idx != -1 {
			raw = raw[idx+1:]
		}

		imageBytes, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			imageBytes, err = base64.RawStdEncoding.DecodeString(raw)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "avatar bukan base64 yang valid")
			}
		}

		mimeType, ext, err := detectContentType(imageBytes)
		if err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("avatars/%s_%d%s",
			req.UserId,
			time.Now().UnixMilli(),
			ext,
		)

		resp, err := r.Storage.PostStorage(c, &v1.PostStorageRequest{
			FileData:    imageBytes,
			FileName:    fileName,
			BucketName:  pkg.BucketProfile,
			ContentType: mimeType,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "upload avatar gagal: %v", err)
		}
		avatarURL = &resp.Url
	}

	query := `
		UPDATE users SET
			avatar_url = COALESCE($1, avatar_url),
			name       = COALESCE(NULLIF($2, ''), name)
		WHERE id = $3
		RETURNING id, name, avatar_url
	`

	var user v1.User
	err := r.DB.DB.QueryRowContext(c, query, avatarURL, req.Name, req.UserId).
		Scan(&user.Id, &user.Name, &user.AvatarUrl)
	if err != nil {
		log.Printf("[update profile] error: %s", err.Error())
		return nil, status.Error(codes.Internal, "failed to update profile")
	}

	return &v1.UpdateProfileResponse{User: &user}, nil
}
