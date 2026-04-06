package userservices

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo UserRepository
}

func NewServices(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}

	return s.repo.GetProfile(ctx, req)
}
func (s *Service) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("username is required")
	}
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}
	if err := validateMaxLength("username", req.Username, 50); err != nil {
		return nil, err
	}
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := validateMaxLength("email", req.Email, 100); err != nil {
		return nil, err
	}

	return s.repo.UpdateUser(ctx, req)
}

func (s *Service) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.OldPassword) == "" {
		return nil, errors.New("old_password is required")
	}
	if err := validatePassword(req.NewPassword); err != nil {
		return nil, err
	}
	if err := validateMaxLength("new_password", req.NewPassword, 20); err != nil { // bcrypt max 72 bytes
		return nil, err
	}
	if req.OldPassword == req.NewPassword {
		return nil, errors.New("new password must be different from old password")
	}

	return s.repo.ChangePassword(ctx, req)
}

func (s *Service) GetUserSessions(ctx context.Context, req *GetUserSessionsRequest) (*GetUserSessionsResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}

	return s.repo.GetUserSessions(ctx, req)
}

func (s *Service) RevokeSession(ctx context.Context, req *RevokeSessionRequest) (*RevokeSessionResponse, error) {
	if strings.TrimSpace(req.SessionId) == "" {
		return nil, errors.New("session_id is required")
	}

	return s.repo.RevokeSession(ctx, req)
}
