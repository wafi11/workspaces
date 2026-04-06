package workspaceservice

import (
	"context"
	"errors"
	"time"
)

const (
	workspaceCacheKey  = "workspace:%s"
	workspacesCacheKey = "workspaces:user:%s"
	cacheTTL           = 5 * time.Minute
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrQuotaExceeded     = errors.New("workspace quota exceeded")
	ErrValidation        = errors.New("validation error")
	ErrTemplateNotFound  = errors.New("template not found")
)

type WorkspaceStatus string

const (
	StatusPending  WorkspaceStatus = "pending"
	StatusRunning  WorkspaceStatus = "running"
	StatusStopped  WorkspaceStatus = "stopped"
	StatusError    WorkspaceStatus = "error"
	StatusDeleting WorkspaceStatus = "deleting"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error)
	ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error)
	ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error)
	GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error)
	UpdateWorkspaceStatus(ctx context.Context, workspaceId string, status WorkspaceStatus) error
}

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error)
	ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error)
	ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error)
	GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error)
}

type Workspace struct {
	Id           string          `json:"id"`
	UserId       string          `json:"user_id"`
	TemplateId   string          `json:"template_id"`
	TemplateName string          `json:"-"`
	Name         string          `json:"name"`
	Namespace    string          `json:"namespace,omitempty"`
	Status       WorkspaceStatus `json:"status"`
	EnvVars      map[string]any  `json:"env_vars"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Url          string          `json:"url"`
}

type CachedWorkspace struct {
	Id         string          `json:"id"`
	UserId     string          `json:"user_id"`
	TemplateId string          `json:"template_id"`
	Name       string          `json:"name"`
	Namespace  string          `json:"namespace"`
	Status     WorkspaceStatus `json:"status"`
	EnvVars    map[string]any  `json:"env_vars"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type CreateWorkspaceRequest struct {
	UserId     string         `json:"user_id"`
	TemplateId string         `json:"template_id"`
	Name       string         `json:"name"`
	EnvVars    map[string]any `json:"env_vars"`
	Addons     []string       `json:"addons"`
}

type CreateWorkspaceResponse struct {
	Workspace *Workspace `json:"workspace"`
	Message   string     `json:"message"`
}

type ListWorkspacesRequest struct {
	UserId string `json:"user_id"`
}

type ListWorkspacesResponse struct {
	Workspaces []Workspace `json:"workspaces"`
}

type GetWorkspaceRequest struct {
	WorkspaceId string `json:"workspace_id"`
	UserId      string `json:"user_id"`
}

type GetWorkspaceResponse struct {
	Workspace *Workspace `json:"workspace"`
}

type DeleteWorkspaceRequest struct {
	WorkspaceId string `json:"workspace_id"`
	UserId      string `json:"user_id"`
}

type DeleteWorkspaceResponse struct {
	Message string `json:"message"`
}
