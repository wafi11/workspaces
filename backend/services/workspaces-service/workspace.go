package workspaceservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wafi11/workspaces/config"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/proto"
	"github.com/wafi11/workspaces/pkg/utils"
	"github.com/wafi11/workspaces/pkg/websocket"
	authservices "github.com/wafi11/workspaces/services/auth-service"
)

type Repository struct {
	db          *sqlx.DB
	redis      *config.RedisConnection
	hub      *websocket.Hub

}

func NewRepository(db *sqlx.DB, redis *config.RedisConnection,	hub      *websocket.Hub) *Repository {
	return &Repository{
		db:          db,
		redis: redis,
		hub: hub,
	}
}

func (r *Repository) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var maxWs, usedWs int

	// Ambil Quota User
	err = tx.QueryRowContext(ctx, `
        SELECT max_workspaces, used_workspaces 
        FROM user_quotas WHERE user_id = $1 FOR UPDATE`,
		req.UserId,
	).Scan(&maxWs, &usedWs)

	// Validasi Quota
	if err != nil {
		return nil, fmt.Errorf("gagal ambil quota: %w", err)
	}

	if usedWs+1 > maxWs {
		return nil, fmt.Errorf("quota workspaces penuh")
	}

	// 1. Insert ke tabel Workspaces
	var w Workspace
	var template_name string

	// if req.Password == "" {
	// 	return nil, fmt.Errorf("password must be required")
	// }
	// hashedPassword, err := utils.HashPassword(req.Password)


	envJSON, _ := json.Marshal(req.EnvVars)

	db_name := utils.GetEnvString(req.EnvVars,"DB_NAME")
	db_password := utils.GetEnvString(req.EnvVars,"DB_PASSWORD")
	db_user := utils.GetEnvString(req.EnvVars,"DB_USER")

	err = tx.QueryRowContext(ctx,`select name from templates where id = $1`,req.TemplateId).Scan(&template_name)

	if err != nil {
		return nil,fmt.Errorf("failed to create workspaces")
	}

	url := GenerateAddonConnectionUrl(AddonUrl(strings.ToLower(template_name)),authservices.GenerateNamespace(req.UserId),req.Name,req.UserId,db_user,db_password,db_name)

	err = tx.QueryRowContext(ctx, `
        INSERT INTO workspaces (user_id,name,status,env_vars,url,template_id)
        VALUES ($1, $2, $3, $4,$5,$6)
        RETURNING id, user_id, name, status`,
		req.UserId, req.Name, StatusPending, envJSON, url,req.TemplateId,
	).Scan(&w.Id, &w.UserId, &w.Name, &w.Status)

	if err != nil {
		log.Printf("failed to create workspaces : %s", err.Error())
		return nil, fmt.Errorf("failed to create workspace")
	}

	var wsResourcesID string
	err = tx.QueryRowContext(ctx,
		`insert into workspace_resources (
			workspace_id,
			name,
			status,
			cpu_cores_req,
			ram_mb_req,
			limit_ram_mb,
			limit_cpu_cores
		) values (
			$1,$2,$3,$4,$5,$6,$7
		) RETURNING id
		`, w.Id, w.Name, StatusPending, req.ReqCpuCores, req.ReqRam, req.LimitRam, req.LimitCpuCores,
	).Scan(&wsResourcesID)

	if err != nil {
		log.Printf("failed to create workspaces : %s", err.Error())
		return nil, fmt.Errorf("failed to create workspace")
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE user_quotas 
       	SET used_workspaces = used_workspaces + 1
		WHERE user_id = $1`,
		req.UserId,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal update quota: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if req.TemplateId != "" {
		messagebroker.PublishEvent(ctx, r.redis.Redis, &proto.WorkspaceEnvelope{
			Payload: &proto.WorkspaceEnvelope_Add{
				Add: &proto.AddPodEvent{
					Identity: &proto.WorkspaceIdentity{
						WorkspaceId: w.Id,
						UserId:      w.UserId,
						Username:    username,
						Name:        w.Name,
						Password:    req.Password,
						Namespace:   w.UserId,
					},
					AddOns:  &proto.AddonSpec{
						Image: "",
						DbUser: db_user,
						DbPassword: db_password,
						DbName: db_name,
					},
					TemplateId: req.TemplateId,
					Resources: &proto.ResourceSpec{
						CpuRequest:    fmt.Sprintf("%.2f", req.ReqCpuCores),
						CpuLimit:      fmt.Sprintf("%.2f", req.LimitCpuCores),
						MemoryRequest: fmt.Sprintf("%dMi", req.ReqRam),
						MemoryLimit:   fmt.Sprintf("%dMi", req.LimitRam),
					},
					Replicas: 1,
				},
			},
		})
	}

	return &CreateWorkspaceResponse{Workspace: &w, Message: "provisioning in progress"}, nil
}

func (r *Repository) ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			w.id, 
			w.url, 
			w.name, 
			w.status, 
			w.env_vars, 
			w.created_at, 
			w.updated_at,
			t.icon,
			ws.started_at,
			ws.stopped_at,
			ws.expires_at,
			ws.next_start_at,
			ws.timezone
		FROM workspaces w
		LEFT JOIN templates t on t.id = w.template_id
		LEFT JOIN LATERAL (
			SELECT started_at, stopped_at, expires_at,next_start_at, timezone
			FROM workspace_sessions
			WHERE workspace_id = w.id
			ORDER BY started_at DESC
			LIMIT 1
		) ws ON true
		WHERE w.user_id = $1 AND w.status != $2
		ORDER BY w.created_at DESC
	`, req.UserId, StatusDeleting)

	if err != nil {
		log.Printf("workspace by user error : %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var workspaces []WorkspaceAndSessions
	for rows.Next() {
		var w WorkspaceAndSessions
		var envRaw []byte
		if err := rows.Scan(
			&w.Id,
			&w.Url,
			&w.Name,
			&w.Status,
			&envRaw,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.Icon,
			&w.StartedAt,
			&w.StoppedAt,
			&w.ExpiresAt,
			&w.NextStartAt,
			&w.Timezone,
		); err != nil {
			log.Printf("workspace by user error : %s", err.Error())
			return nil, err
		}
		if err := json.Unmarshal(envRaw, &w.EnvVars); err != nil {
			w.EnvVars = map[string]any{}
		}
		workspaces = append(workspaces, w)
	}

	return &ListWorkspacesResponse{Workspaces: workspaces}, nil
}

func (r *Repository) ListWorkspaceForm(ctx context.Context, userId string) ([]ListWorkspaceForm, error) {
	query := `
		select id,name from workspaces where user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		log.Printf("failed to scan list Workspace : %s \n", err.Error())
		return nil, fmt.Errorf("workspace not found")
	}
	defer rows.Close()

	var datas []ListWorkspaceForm
	for rows.Next() {
		var data ListWorkspaceForm
		err = rows.Scan(&data.Id, &data.Name)
		if err != nil {
			log.Printf("failed to scan list Workspace : %s \n", err.Error())
			return nil, fmt.Errorf("workspace not found")
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (r *Repository) ListWorkspaces(ctx context.Context, limit int, offset int, status string) (*ListWorkspacesResponse, error) {
	if status == "" {
		status = "running"
	}

	rows, err := r.db.QueryContext(ctx, `
      SELECT 
			w.id, 
			w.url, 
			w.name, 
			w.status, 
			w.env_vars, 
			w.created_at, 
			w.updated_at,
			t.icon,
			ws.started_at,
			ws.stopped_at,
			ws.expires_at,
			ws.timezone
		FROM workspaces w
		LEFT JOIN templates t on t.id = w.template_id
		LEFT JOIN LATERAL (
			SELECT started_at, stopped_at, expires_at, timezone
			FROM workspace_sessions
			WHERE workspace_id = w.id
			ORDER BY started_at DESC
			LIMIT 1
		) ws ON true
		WHERE status = $1
        ORDER BY user_id DESC,created_at DESC
        LIMIT $2 OFFSET $3
    `, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []WorkspaceAndSessions
	for rows.Next() {
		var w WorkspaceAndSessions
		var envRaw []byte
		if err := rows.Scan(
			&w.Id,
			&w.Url,
			&w.Name,
			&w.Status,
			&envRaw,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.Icon,
			&w.StartedAt,
			&w.StoppedAt,
			&w.ExpiresAt,
			&w.Timezone,
		); err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		if err := json.Unmarshal(envRaw, &w.EnvVars); err != nil {
			w.EnvVars = map[string]any{}
		}
		workspaces = append(workspaces, w)
	}

	return &ListWorkspacesResponse{Workspaces: workspaces}, nil
}

func (r *Repository) GetWorkspace(ctx context.Context, req *GetWorkspaceRequest) (*GetWorkspaceResponse, error) {
	// 1. check cache
	if cached, err := r.getWorkspaceCache(ctx, req.WorkspaceId); err == nil {
		return &GetWorkspaceResponse{Workspace: cached}, nil
	}

	// 2. hit db
	var w Workspace
	var envRaw []byte
	err := r.db.QueryRowContext(ctx, `
    SELECT 
        w.name,
        w.status,
        w.url,
        t.name, 
        t.icon, 
        w.env_vars, 
        t.created_at
    FROM workspaces w 
    LEFT JOIN templates t ON t.id = w.template_id
    WHERE w.id = $1
`, req.WorkspaceId).Scan(
    &w.Name,         
    &w.Status,      
    &w.Url,          
    &w.TemplateName, 
    &w.Icon,         
    &envRaw,        
    &w.CreatedAt,   
)
	if err != nil {
		log.Printf("template error : %s",err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal(envRaw, &w.EnvVars); err != nil {
		w.EnvVars = map[string]any{}
	}

	r.setWorkspaceCache(ctx, &w)

	return &GetWorkspaceResponse{Workspace: &w}, nil
}

func (r *Repository) DeleteWorkspace(ctx context.Context, req *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error) {
	// 1. cek ownership
	var w Workspace
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM workspaces WHERE id = $1 AND user_id = $2`,
		req.WorkspaceId, req.UserId,
	).Scan(&w.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, err
	}

	// 2. update status → deleting
	_, err = r.db.ExecContext(ctx,
		`UPDATE workspaces SET status = $1, updated_at = NOW() WHERE id = $2`,
		StatusDeleting, req.WorkspaceId,
	)
	if err != nil {
		return nil, err
	}

	// 4. invalidate cache
	r.invalidateWorkspaceCache(ctx, req.WorkspaceId, req.UserId)

	return &DeleteWorkspaceResponse{Message: "workspace is being deleted"}, nil
}
