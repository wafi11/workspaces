package internal

import (
	"context"

	v1 "github.com/wafi11/workspaces/core/user-service/gen/v1"
	"github.com/wafi11/workspaces/core/user-service/pkg"
)

type Service struct {
	repo pkg.IRepository
	v1.UnimplementedUserServiceServer
}

func NewService(repo pkg.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProfile(c context.Context, req *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	return s.repo.GetProfile(c, req)
}

func (s *Service) UpdateProfile(c context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {

	return s.repo.UpdateProfile(c, req)
}
