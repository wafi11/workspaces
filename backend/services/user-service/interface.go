package userservices

import (
	"context"
	"time"
)

const (
	userCacheKey    = "user:%s"
	sessionCacheKey = "sessions:%s"
	cacheTTL        = 5 * time.Minute
)

type UserRepository interface {
	GetProfile(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error)
	GetUserSessions(ctx context.Context, req *GetUserSessionsRequest) (*GetUserSessionsResponse, error)
	RevokeSession(ctx context.Context, req *RevokeSessionRequest) (*RevokeSessionResponse, error)
}

type UserService interface {
	GetProfile(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error)
	GetUserSessions(ctx context.Context, req *GetUserSessionsRequest) (*GetUserSessionsResponse, error)
	RevokeSession(ctx context.Context, req *RevokeSessionRequest) (*RevokeSessionResponse, error)
}

type CachedUser struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CachedSession struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserSessionsRequest struct {
	UserId string `json:"user_id"`
}

type GetUserSessionsResponse struct {
	Sessions []Session `json:"sessions"`
}

type RevokeSessionRequest struct {
	SessionId string `json:"session_id"`
}

type RevokeSessionResponse struct {
	Message string `json:"message"`
}

type GetUserRequest struct {
	UserId string `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type UpdateUserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}

type ChangePasswordRequest struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
