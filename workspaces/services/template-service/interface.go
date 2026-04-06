package templateservice

import (
	"context"
	"errors"
	"time"
)

type TemplateRepository interface {
	ExecuteDeployment(ctx context.Context, templateName string, params DeployParams) error
	ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error)
	GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error)
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id string, template *UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id string) error
}

type TemplateService interface {
	ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error)
	GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error)
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id string, template *UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id string) error
}

const (
	templateCacheKey  = "template:%s"
	templatesCacheKey = "templates:all"
	cacheTTL          = 5 * time.Minute
)

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrValidation       = errors.New("validation error")
)

type Template struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	Icon        string            `json:"icon"`
	Category    string             `json:"category"`
	IsPublic    bool               `json:"is_public"`
	TemplateUrl string            `json:"template_url"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
	Addons      []TemplateAddon    `json:"addons,omitempty"`
	Files       []TemplateFiles    `json:"files,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

type TemplateFiles struct {
	Id         string `json:"id"`
	Filename   string `json:"filename"`
	TemplateId string `json:"template_id"`
	SortOrder  int    `json:"sort_order"`
}

type TemplateVariable struct {
	Id           string `json:"id"`
	TemplateId   string `json:"template_id"`
	Key          string `json:"key"`
	DefaultValue string `json:"default_value"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
}

type TemplateAddon struct {
	Id            string         `json:"id"`
	TemplateId    string         `json:"template_id"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Description   string         `json:"description"`
	DefaultConfig map[string]any `json:"default_config"`
}

type TemplateFileInfo struct {
	TemplateUrl string `db:"template_url"`
	Filename    string `db:"filename"`
	SortOrder   int    `db:"sort_order"`
}

type DeployParams struct {
	DbName       *string
	User         *string
	Name         string
	Namespace    string
	StorageClass string
	StorageSize  string
	Replicas     int
	RunAsUser    int
	RunAsGroup   int
	FsGroup      int
	Password     string
	CPURequest   string
	MemRequest   string
	CPULimit     string
	MemLimit     string
	Username     string
	Domain       string
}

type CachedTemplate struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	Category    string             `json:"category"`
	IsPublic    bool               `json:"is_public"`
	TemplateUrl string            `json:"template_url"`
	Icon        string             `json:"icon"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
	Addons      []TemplateAddon    `json:"addons,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

type ListTemplatesRequest struct {
	Category string `json:"category"`
}

type ListTemplatesResponse struct {
	Templates []Template `json:"templates"`
}

type GetTemplateRequest struct {
	TemplateId string `json:"template_id"`
}

type GetTemplateResponse struct {
	Template *Template `json:"template"`
}

type CreateTemplateRequest struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Image       string                       `json:"image"`
	Category    string                       `json:"category"`
	IsPublic    bool                         `json:"is_public"`
	Icon        string                       `json:"icon"`
	Variables   []CreateVariableRequest      `json:"variables"`
	Addons      []CreateAddonRequest         `json:"addons"`
	Files       []CreateTemplateFilesRequest `json:"files"`
}

type CreateVariableRequest struct {
	Key          string `json:"key"`
	DefaultValue string `json:"default_value"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
}

type CreateAddonRequest struct {
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Description   string         `json:"description"`
	DefaultConfig map[string]any `json:"default_config"`
}

type CreateTemplateFilesRequest struct {
	Filename  string `json:"filename"`
	SortOrder int    `json:"sort_order"`
}

type CreateTemplateResponse struct {
	Template *Template `json:"template"`
	Message  string    `json:"message"`
}

type UpdateTemplateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Category    *string `json:"category"`
	IsPublic    *bool   `json:"is_public"`
	TemplateUrl *string `json:"template_url"`
	Icon        *string `json:"icon"`
}

func (req *UpdateTemplateRequest) merge(t *Template) {
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Image != nil {
		t.Image = *req.Image
	}
	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.IsPublic != nil {
		t.IsPublic = *req.IsPublic
	}
	if req.TemplateUrl != nil {
		t.TemplateUrl = *req.TemplateUrl
	}
	if req.Icon != nil {
		t.Icon = *req.Icon
	}
}
