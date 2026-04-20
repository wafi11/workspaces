package templateservice

import (
	"context"
	"database/sql"

	"github.com/wafi11/workspaces/pkg/models"
)

type TemplateRepository interface {
	ListTemplates(ctx context.Context, req *models.ListTemplatesRequest) (*models.ListTemplatesResponse, error)
	GetTemplate(ctx context.Context, req *models.GetTemplateRequest) (*models.GetTemplateResponse, error)
	CreateTemplate(ctx context.Context, req *models.CreateTemplateRequest) (*models.CreateTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id string, template *models.UpdateTemplateRequest) error
	GetDetailsInfo(c context.Context, templateId string) (*models.DetailsInfo, error)
	DeleteTemplate(ctx context.Context, id string) error
	FindTemplateWorkspaceForm(c context.Context)([]models.TemplateWorkspaceForm,error)

	// template-variables
	CreateTemplateVariable(ctx context.Context, req *models.CreateVariableRequest, templateId string,Tx *sql.Tx) error
	DeleteTemplateVariable(ctx context.Context, id string) error
	GetTemplateVariables(ctx context.Context, templateID string) ([]models.TemplateVariable, error)
	UpdateTemplateVariable(ctx context.Context, id string, req *models.CreateVariableRequest) error

	// template-addons
	CreateTemplateAddon(ctx context.Context, req *models.CreateAddonRequest, templateId string,Tx *sql.Tx) error
	DeleteTemplateAddon(ctx context.Context, id string) error
	GetTemplateAddons(ctx context.Context, templateID string) ([]models.TemplateAddon, error)
	UpdateTemplateAddon(ctx context.Context, id string, req *models.CreateAddonRequest) error

	// template-files
	CreateTemplateFiles(ctx context.Context, req *models.CreateTemplateFilesRequest, templateId string,Tx *sql.Tx) error
	DeleteTemplateFiles(ctx context.Context, id string) error
	GetTemplateFiles(ctx context.Context, templateID string) ([]models.TemplateFiles, error)
	UpdateTemplateFiles(ctx context.Context, id string, req *models.CreateTemplateFilesRequest) error
}

type TemplateService interface {
	ListTemplates(ctx context.Context, req *models.ListTemplatesRequest) (*models.ListTemplatesResponse, error)
	GetTemplate(ctx context.Context, req *models.GetTemplateRequest) (*models.GetTemplateResponse, error)
	CreateTemplate(ctx context.Context, req *models.CreateTemplateRequest) (*models.CreateTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id string, template *models.UpdateTemplateRequest) error
	GetDetailsInfo(c context.Context, templateId string) (*models.DetailsInfo, error)
	DeleteTemplate(ctx context.Context, id string) error
	FindTemplateWorkspaceForm(c context.Context)([]models.TemplateWorkspaceForm,error)

	// template-variables
	CreateTemplateVariable(ctx context.Context, req *models.CreateVariableRequest, templateId string) error
	DeleteTemplateVariable(ctx context.Context, id string) error
	GetTemplateVariables(ctx context.Context, templateID string) ([]models.TemplateVariable, error)
	UpdateTemplateVariable(ctx context.Context, id string, req *models.CreateVariableRequest) error

	// template-addons
	CreateTemplateAddon(ctx context.Context, req *models.CreateAddonRequest, templateId string) error
	DeleteTemplateAddon(ctx context.Context, id string) error
	GetTemplateAddons(ctx context.Context, templateID string) ([]models.TemplateAddon, error)
	UpdateTemplateAddon(ctx context.Context, id string, req *models.CreateAddonRequest) error

	// template-files
	CreateTemplateFiles(ctx context.Context, req *models.CreateTemplateFilesRequest, templateId string) error
	DeleteTemplateFiles(ctx context.Context, id string) error
	GetTemplateFiles(ctx context.Context, templateID string) ([]models.TemplateFiles, error)
	UpdateTemplateFiles(ctx context.Context, id string, req *models.CreateTemplateFilesRequest) error
}
