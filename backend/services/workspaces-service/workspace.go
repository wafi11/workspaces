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
)

func GenerateNamespace(userId, name string) string {
	return fmt.Sprintf("ws-%s-%s", userId[:8], name)
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

// ─── Repository Methods ───────────────────────────────────────────────────────
func (r *Repository) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest, username string) (*CreateWorkspaceResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Lock quota row dulu biar ga bisa diakses concurrent
	var maxWorkspaces, used_workspaces int
	err = tx.QueryRowContext(ctx,
		`SELECT used_workspaces,max_workspaces FROM user_quotas WHERE user_id = $1 FOR UPDATE`,
		req.UserId,
	).Scan(&maxWorkspaces, &used_workspaces)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			maxWorkspaces = 3
		} else {
			return nil, err
		}
	}

	if used_workspaces >= maxWorkspaces {
		return nil, ErrQuotaExceeded
	}

	// 3. Insert workspace
	var w Workspace
	var envJSON, envRaw []byte
	envJSON, err = json.Marshal(req.EnvVars)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, `
        INSERT INTO workspaces (user_id, template_id, name, namespace, status, env_vars)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, user_id, template_id, name, namespace, status, env_vars, created_at, updated_at
    `, req.UserId, req.TemplateId, req.Name, GenerateNamespace(req.UserId, req.Name), StatusPending, envJSON,
	).Scan(
		&w.Id, &w.UserId, &w.TemplateId,
		&w.Name, &w.Namespace, &w.Status,
		&envRaw, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(envRaw, &w.EnvVars); err != nil {
		w.EnvVars = map[string]any{}
	}

	// 4. Update URL dalam transaction yang sama
	w.Url = generateUrl(w.Id)
	_, err = tx.ExecContext(ctx, `UPDATE workspaces SET url = $1 WHERE id = $2`, w.Url, w.Id)
	if err != nil {
		return nil, err
	}
	var q struct {
		MaxWS, MaxCPU, MaxRAM, MaxStorage     int
		UsedWS, UsedCPU, UsedRAM, UsedStorage int
	}
	err = tx.QueryRowContext(ctx, `
        SELECT 
            max_workspaces, max_cpu_cores, max_ram_mb, max_storage_gb,
            used_workspaces, used_cpu_cores, used_ram_mb, used_storage_gb
        FROM user_quotas WHERE user_id = $1 FOR UPDATE`,
		req.UserId,
	).Scan(
		&q.MaxWS, &q.MaxCPU, &q.MaxRAM, &q.MaxStorage,
		&q.UsedWS, &q.UsedCPU, &q.UsedRAM, &q.UsedStorage,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Default quota jika belum ada record
			q.MaxWS, q.MaxCPU, q.MaxRAM, q.MaxStorage = 3, 4, 4096, 20
		} else {
			return nil, err
		}
	}
	reqCPU := 1
	reqRAM := 1024
	reqSTG := 5

	// 3. Validasi Quota
	if q.UsedWS+1 > q.MaxWS {
		return nil, fmt.Errorf("limit workspace tercapai")
	}
	if q.UsedCPU+reqCPU > q.MaxCPU {
		return nil, fmt.Errorf("limit CPU tidak mencukupi")
	}
	if q.UsedRAM+reqRAM > q.MaxRAM {
		return nil, fmt.Errorf("limit RAM tidak mencukupi")
	}
	_, err = tx.ExecContext(ctx, `
        INSERT INTO workspace_resources (workspace_id, kind, name, cpu_cores, ram_mb, storage_gb, status)
        VALUES ($1, 'deployment', $2, $3, $4, $5, 'pending')`,
		w.Id, w.Name, reqCPU, reqRAM, reqSTG,
	)
	if err != nil {
		return nil, err
	}

	// 5. Commit dulu baru publish event
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 6. Publish event & invalidate cache (di luar transaction)
	PublishEvent(ctx, r.redisClient, WorkspaceJob{
		WorkspaceId:            w.Id,
		UserId:                 w.UserId,
		Namespace:              w.Namespace,
		TemplateId:             req.TemplateId,
		Username:               username,
		Action:                 JobCreate,
		EnvVars:                w.EnvVars,
		MemoryRequest:          "100Mi",
		StorageRequest:         "1Gi",
		MemoryTerminalRequest:  "",
		StorageTerminalRequest: "",
		CpuTerminalReuqest:     "",
		Replicas:               "1",
		CPURequest:             "0.25",
	})

	r.redisClient.Del(ctx, fmt.Sprintf(workspacesCacheKey, req.UserId))

	return &CreateWorkspaceResponse{
		Workspace: &w,
		Message:   "workspace created, provisioning in progress",
	}, nil
}

func (r *Repository) ListWorkspacesByUserId(ctx context.Context, req *ListWorkspacesRequest) (*ListWorkspacesResponse, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, template_id, name, status, env_vars, created_at, updated_at
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
			&w.Id, &w.UserId, &w.TemplateId,
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
        SELECT id, user_id, template_id, name, namespace, status, env_vars, created_at, updated_at
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
			&w.Id, &w.UserId, &w.TemplateId,
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
		SELECT id, user_id, template_id, name,namespace, status,url, env_vars, created_at, updated_at
		FROM workspaces
		WHERE id = $1
	`, req.WorkspaceId,
	).Scan(
		&w.Id, &w.UserId, &w.TemplateId,
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

	PublishEvent(ctx, r.redisClient, WorkspaceJob{
		WorkspaceId: w.Id,
		UserId:      w.UserId,
		Namespace:   w.Namespace,
		Action:      JobDelete,
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
