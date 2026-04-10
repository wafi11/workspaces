package workspaceservice

import (
	"context"
	"errors"
	"log"
)

type Service struct {
	repo WorkspaceRepository
}

func NewService(repo WorkspaceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error) {
	if err := validateCreateWorkspace(req); err != nil {
		return nil, err
	}
	return s.repo.CreateWorkspace(ctx, req, username)
}

func (s *Service) ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error) {
	if err := validateListWorkspaces(req); err != nil {
		return nil, err
	}
	return s.repo.ListWorkspacesByUserId(ctx, req)
}

func (s *Service) ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error) {
	data, err := s.repo.ListWorkspaces(ctx, limit, offset, status)
	if err != nil {
		log.Printf("ERROR ListWorkspaces: %v", err)
		return nil, errors.New("failed to fetch workspaces")
	}
	return &ListWorkspacesResponse{
		Workspaces: data.Workspaces,
	}, nil
}

func (s *Service) ListWorkspaceForm(ctx context.Context, userId string) ([]ListWorkspaceForm, error) {
	return s.repo.ListWorkspaceForm(ctx, userId)
}

func (s *Service) GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error) {
	if err := validateGetWorkspace(req, true); err != nil {
		return nil, err
	}
	return s.repo.GetWorkspace(ctx, req)
}

func (s *Service) DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error) {
	if err := validateDeleteWorkspace(req); err != nil {
		return nil, err
	}
	return s.repo.DeleteWorkspace(ctx, req)
}

// workspace addon
func (s *Service) GetAddonService(c context.Context, workspaceId string) ([]WorkspaceAddonDetails, error) {
	return s.repo.GetAddonService(c, workspaceId)
}

func (s *Service) CreateAddonWorkspace(c context.Context, req CreateWorkspaceAddon) error {
	return s.repo.CreateAddonWorkspace(c, req)
}

func (s *Service) UpdateWorkspaceStatus(ctx context.Context, workspaceId string, status WorkspaceStatus) error {
	return s.repo.UpdateWorkspaceStatus(ctx, workspaceId, status)
}
