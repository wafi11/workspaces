package templateservice

import "context"

type Service struct {
	repo TemplateRepository
}

func NewService(repo TemplateRepository) *Service {
	return &Service{repo: repo}
}

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

func (s *Service) UpdateTemplate(ctx context.Context,id string, req *UpdateTemplateRequest) error {
	return s.repo.UpdateTemplate(ctx,id,req)
}

func (s *Service) DeleteTemplate(ctx context.Context,id string) error {
	return s.repo.DeleteTemplate(ctx,id)
}
