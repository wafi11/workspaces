package services

import (
	"errors"
	"fmt"
	"time"
)

func generateNamespace(userID string) string {
	return fmt.Sprintf("ws-%s", userID)
}

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrQuotaExceeded     = errors.New("workspace quota exceeded")
	ErrValidation        = errors.New("validation error")
	ErrTemplateNotFound  = errors.New("template not found")
)

const (
	workspaceCacheKey  = "workspace:%s"
	workspacesCacheKey = "workspaces:user:%s"
	cacheTTL           = 5 * time.Minute
)

type WorkspaceStatus string

const (
	StatusPending  WorkspaceStatus = "pending"
	StatusRunning  WorkspaceStatus = "running"
	StatusStopped  WorkspaceStatus = "stopped"
	StatusError    WorkspaceStatus = "error"
	StatusDeleting WorkspaceStatus = "deleting"
)

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

type DeployParams struct {
	WsID             string
	WS_TOKEN         string
	WS_REFRESH_TOKEN string
	WS_API_URL       string
	DB_USER          *string
	DB_NAME          *string
	DB_PASSWORD      *string
	Image            *string
	User             *string
	Name             string
	Namespace        string
	StorageClass     string
	StorageSize      string
	Replicas         int
	RunAsUser        int
	RunAsGroup       int
	FsGroup          int
	Password         string
	CPURequest       string
	MemRequest       string
	CPULimit         string
	MemLimit         string
	Username         string
	Domain           string
}

type TemplateFileInfo struct {
	TemplateUrl string `db:"template_url"`
	Filename    string `db:"filename"`
	SortOrder   int    `db:"sort_order"`
}
