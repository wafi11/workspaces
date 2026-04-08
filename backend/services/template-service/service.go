package templateservice

import (
	"context"
	"fmt"
)

type Service struct {
	repo TemplateRepository
}

func NewService(repo TemplateRepository) *Service {
	return &Service{repo: repo}
}

// templates
func (s *Service) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	if err := validateListTemplates(req); err != nil {
		return nil, err
	}
	return s.repo.ListTemplates(ctx, req)
}

func (s *Service) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {

	return s.repo.GetTemplate(ctx, req)
}

func (s *Service) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	if err := validateCreateTemplate(req); err != nil {
		return nil, err
	}
	return s.repo.CreateTemplate(ctx, req)
}

func (s *Service) UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) error {
	return s.repo.UpdateTemplate(ctx, id, req)
}

func (s *Service) GetDetailsInfo(c context.Context, templateId string) (*DetailsInfo, error) {
	if templateId == "" {
		return nil, fmt.Errorf("templates not found")
	}
	return s.repo.GetDetailsInfo(c, templateId)
}

func (s *Service) DeleteTemplate(ctx context.Context, id string) error {
	return s.repo.DeleteTemplate(ctx, id)
}

// template-variables
func (s *Service) CreateTemplateVariable(ctx context.Context, req *CreateVariableRequest, templateId string) error {
	return s.repo.CreateTemplateVariable(ctx, req, templateId)
}
func (s *Service) DeleteTemplateVariable(ctx context.Context, id string) error {
	return s.repo.DeleteTemplateVariable(ctx, id)
}
func (s *Service) GetTemplateVariables(ctx context.Context, templateID string) ([]TemplateVariable, error) {
	return s.repo.GetTemplateVariables(ctx, templateID)
}
func (s *Service) UpdateTemplateVariable(ctx context.Context, id string, req *CreateVariableRequest) error {
	return s.repo.UpdateTemplateVariable(ctx, id, req)
}

// template-addons
func (s *Service) CreateTemplateAddon(ctx context.Context, req *CreateAddonRequest, templateId string) error {
	return s.repo.CreateTemplateAddon(ctx, req, templateId)
}
func (s *Service) DeleteTemplateAddon(ctx context.Context, id string) error {
	return s.repo.DeleteTemplateAddon(ctx, id)
}
func (s *Service) GetTemplateAddons(ctx context.Context, templateID string) ([]TemplateAddon, error) {
	return s.repo.GetTemplateAddons(ctx, templateID)
}
func (s *Service) UpdateTemplateAddon(ctx context.Context, id string, req *CreateAddonRequest) error {
	return s.repo.UpdateTemplateAddon(ctx, id, req)
}

// template-files
func (s *Service) CreateTemplateFiles(ctx context.Context, req *CreateTemplateFilesRequest, templateId string) error {
	return s.repo.CreateTemplateFiles(ctx, req, templateId)
}

func (s *Service) DeleteTemplateFiles(ctx context.Context, id string) error {
	return s.repo.DeleteTemplateFiles(ctx, id)
}

func (s *Service) GetTemplateFiles(ctx context.Context, templateID string) ([]TemplateFiles, error) {
	return s.repo.GetTemplateFiles(ctx, templateID)
}

func (s *Service) UpdateTemplateFiles(ctx context.Context, id string, req *CreateTemplateFilesRequest) error {
	return s.repo.UpdateTemplateFiles(ctx, id, req)
}
