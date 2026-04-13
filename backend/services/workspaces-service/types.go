package workspaceservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wafi11/workspaces/pkg/proto"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error)
	ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error)
	ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error)
	ListWorkspaceForm(ctx context.Context, userId string) ([]ListWorkspaceForm, error)
	GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error)
	UpdateWorkspaceStatus(ctx context.Context, workspaceId string, status string) error
	CreateWorkspaceSessions(ctx context.Context, req CreateWorkspaceSessions) error
	CanStartWorkspace(ctx context.Context, workspaceID string,tx *sql.Tx) (bool, error)
	AutoStopWorkspace(ctx context.Context, workspaceId string) error 
}

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error)
	ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error)
	ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error)
	ListWorkspaceForm(ctx context.Context, userId string) ([]ListWorkspaceForm, error)
	GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error)
	UpdateWorkspaceStatus(ctx context.Context, workspaceId,userId string, status string) error
	DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error)
	StartEventConsumer(ctx context.Context)
}

const (
	workspaceCacheKey  = "workspace:%s"
	workspacesCacheKey = "workspaces:user:%s"
	cacheTTL           = 5 * time.Minute
	cooldown = 2
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrQuotaExceeded     = errors.New("workspace quota exceeded")
	ErrValidation        = errors.New("validation error")
	ErrTemplateNotFound  = errors.New("template not found")
	ErrAddonNotFound     = fmt.Errorf("addon not found")
)

type WorkspaceStatus string

const (
	StatusPending  WorkspaceStatus = "pending"
	StatusRunning  WorkspaceStatus = "running"
	StatusStopped  WorkspaceStatus = "stopped"
	StatusError    WorkspaceStatus = "error"
	StatusDeleting WorkspaceStatus = "deleting"
)

type AddonUrl string 

const (
	PostgresqlURL  AddonUrl = "postgres"
	MysqlURL  AddonUrl = "mysql"
	RedisURL  AddonUrl = "redis"
)

func ConvertWorkspaceStatus(s proto.WorkspaceStatus) WorkspaceStatus {
	switch s {
	case proto.WorkspaceStatus_WORKSPACE_STATUS_PROVISIONING:
		return StatusPending
	case proto.WorkspaceStatus_WORKSPACE_STATUS_RUNNING:
		return StatusRunning
	case proto.WorkspaceStatus_WORKSPACE_STATUS_STOPPED:
		return StatusStopped
	case proto.WorkspaceStatus_WORKSPACE_STATUS_FAILED:
		return StatusError
	case proto.WorkspaceStatus_WORKSPACE_STATUS_DELETING:
		return StatusDeleting
	default:
		return StatusPending
	}
}

type Workspace struct {
	Id        string          `json:"id"`
	UserId    string          `json:"user_id"`
	Name      string          `json:"name"`
	Namespace string          `json:"namespace,omitempty"`
	Status    WorkspaceStatus `json:"status"`
	Icon      *string         `json:"icon"`
	EnvVars   map[string]any  `json:"env_vars"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Url       string          `json:"url"`
}

type WorkspaceAndSessions struct {
	Workspace
	StartedAt *time.Time `json:"started_at"`
	StoppedAt *time.Time `json:"stopped_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	NextStartAt *time.Time `json:"next_start_at"`
	Timezone  *string    `json:"timezone"`
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

type WorkspaceAddon struct {
	ID              string          `db:"id"              json:"id"`
	WorkspaceID     string          `db:"workspace_id"    json:"workspace_id"`
	TemplateAddonId string          `db:"template_addon_id" json:"template_addon_id"`
	Status          string          `db:"status"          json:"status"`
	Config          json.RawMessage `db:"config"          json:"config"`
}

type WorkspaceAddonDetails struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Icon   string          `json:"icon"`
	Status string          `json:"status"`
	Config json.RawMessage `json:"config"`
}

type AddonConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateWorkspaceAddon struct {
	WorkspaceID     string        `json:"workspace_id"`
	TemplateAddonId string        `json:"template_addon_id"`
	Status          string        `json:"status"`
	Config          []AddonConfig `json:"config"`
}

type CreateWorkspaceRequest struct {
	UserId        string            `json:"user_id"`
	TemplateId    string            `json:"template_id"`
	Password      string            `json:"password"`
	Name          string            `json:"name"`
	LimitRam      int               `json:"limit_ram_mb"`
	LimitCpuCores float64           `json:"limit_cpu_cores"`
	ReqRam        int               `json:"req_ram_mb"`
	ReqCpuCores   float64           `json:"req_cpu_cores"`
	EnvVars       map[string]string `json:"env_vars"`
}

type CreateWorkspaceResponse struct {
	Workspace *Workspace `json:"workspace"`
	Message   string     `json:"message"`
}

type ListWorkspacesRequest struct {
	UserId string `json:"user_id"`
}

type ListWorkspacesResponse struct {
	Workspaces []WorkspaceAndSessions `json:"workspaces"`
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

type ListWorkspaceForm struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateWorkspaceSessions struct {
	WorkspaceId string `json:"workspace_id"`
	UserId      string `json:"user_id"`
	Status      string `json:"status"`
	IpAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
}

type workspaceRow struct {
    UserId     string
    Name       string
    CurrStatus WorkspaceStatus
    LimitRAM   int
    LimitCPU   float64
    ReqRAM     int
    ReqCPU     float64
}
