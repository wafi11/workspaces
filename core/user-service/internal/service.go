package internal

import (
	"context"
	"log"

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
	log.Printf("request incoming")
	return s.repo.GetProfile(c, req)
}
