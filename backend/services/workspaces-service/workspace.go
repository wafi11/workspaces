package workspaceservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/pkg/publisher"
)

func GenerateNamespace(userId string) string {
	return fmt.Sprintf("ws-%s", userId)
}

type Repository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		db:          db,
		redisClient: redis,
	}
}

func (r *Repository) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var q struct {
		MaxWS, MaxRAM, MaxStorage    int
		UsedWS, UsedRAM, UsedStorage int
		MaxCPU, UsedCPU              float64
	}

	// Ambil Quota User
	err = tx.QueryRowContext(ctx, `
        SELECT max_workspaces, max_cpu_cores, max_ram_mb, max_storage_gb,
               used_workspaces, used_cpu_cores, used_ram_mb, used_storage_gb
        FROM user_quotas WHERE user_id = $1 FOR UPDATE`,
		req.UserId,
	).Scan(
		&q.MaxWS, &q.MaxCPU, &q.MaxRAM, &q.MaxStorage,
		&q.UsedWS, &q.UsedCPU, &q.UsedRAM, &q.UsedStorage,
	)

	// Resource constants (Bisa ditaruh di config)
	const (
		reqCPU, reqRAM   = 0.5, 512  // Resource untuk Terminal
		codeCpu, codeRAM = 1.0, 1024 // Resource untuk VS Code
		// Storage tidak ditambah lagi kalau sudah ada PVC di namespace ini
	)

	totalReqCPU := reqCPU + codeCpu
	totalReqRAM := reqRAM + codeRAM

	// Validasi Quota
	if q.UsedWS+1 > q.MaxWS {
		return nil, fmt.Errorf("quota workspaces penuh")
	}
	if q.UsedCPU+totalReqCPU > q.MaxCPU {
		return nil, fmt.Errorf("quota CPU tidak cukup")
	}

	// 1. Insert ke tabel Workspaces
	var w Workspace
	envJSON, _ := json.Marshal(req.EnvVars)
	err = tx.QueryRowContext(ctx, `
        INSERT INTO workspaces (user_id, name, status, env_vars)
        VALUES ($1, $2, 'pending', $3)
        RETURNING id, user_id, name, status`,
		req.UserId, req.Name, envJSON,
	).Scan(&w.Id, &w.UserId, &w.Name, &w.Status)

	if err != nil {
		return nil, fmt.Errorf("errr : %s", err.Error())
	}
	// 2. Update Kuota (Hanya CPU & RAM, Storage tetap karena sharing)
	_, err = tx.ExecContext(ctx, `
        UPDATE user_quotas 
        SET used_workspaces = used_workspaces + 1,
            used_cpu_cores = used_cpu_cores + $1,
            used_ram_mb = used_ram_mb + $2
        WHERE user_id = $3`,
		totalReqCPU, totalReqRAM, req.UserId,
	)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	publisher.PublishEvent(ctx, r.redisClient, publisher.WorkspaceJob{
		UserId:        w.UserId,
		TemplateId:    req.TemplateId,
		Name:          "code",
		Username:      username,
		Action:        publisher.JobAdd,
		EnvVars:       req.EnvVars,
		CPURequest:    fmt.Sprintf("%.2f", reqCPU),
		MemoryRequest: fmt.Sprintf("%dMi", reqRAM),
		CPULimit:      fmt.Sprintf("%.2f", codeCpu),
		MemoryLimit:   fmt.Sprintf("%dMi", codeRAM),
	})

	return &CreateWorkspaceResponse{Workspace: &w, Message: "provisioning in progress"}, nil
}

func (r *Repository) ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, status, env_vars, created_at, updated_at
		FROM workspaces
		WHERE user_id = $1 AND status != $2
		ORDER BY created_at DESC
	`, req.UserId, StatusDeleting)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var w Workspace
		var envRaw []byte
		if err := rows.Scan(
			&w.Id, &w.UserId,
			&w.Name, &w.Status,
			&envRaw, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
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
        SELECT id, user_id, name, namespace, status, env_vars, created_at, updated_at
        FROM workspaces
        WHERE status = $1
        ORDER BY user_id DESC,created_at DESC
        LIMIT $2 OFFSET $3
    `, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var w Workspace
		var envRaw []byte
		if err := rows.Scan(
			&w.Id, &w.UserId,
			&w.Name, &w.Namespace, &w.Status,
			&envRaw, &w.CreatedAt, &w.UpdatedAt,
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
		SELECT id, user_id, name,namespace, status,url, env_vars, created_at, updated_at
		FROM workspaces
		WHERE id = $1
	`, req.WorkspaceId,
	).Scan(
		&w.Id, &w.UserId,
		&w.Name, &w.Namespace, &w.Status, &w.Url,
		&envRaw, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
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
		`SELECT id, namespace FROM workspaces WHERE id = $1 AND user_id = $2`,
		req.WorkspaceId, req.UserId,
	).Scan(&w.Id, &w.Namespace)
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

	publisher.PublishEvent(ctx, r.redisClient, publisher.WorkspaceJob{
		WorkspaceId: w.Id,
		UserId:      w.UserId,
		Namespace:   w.Namespace,
		Action:      publisher.JobDelete,
	})

	// 4. invalidate cache
	r.invalidateWorkspaceCache(ctx, req.WorkspaceId, req.UserId)

	return &DeleteWorkspaceResponse{Message: "workspace is being deleted"}, nil
}

func (r *Repository) UpdateWorkspaceStatus(ctx context.Context, workspaceId string, status WorkspaceStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE workspaces SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, workspaceId,
	)
	return err
}
