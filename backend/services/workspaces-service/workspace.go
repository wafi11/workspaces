package workspaceservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
	// 1. cek quota user
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM workspaces WHERE user_id = $1 AND status != $2`,
		req.UserId, StatusDeleting,
	).Scan(&count)
	if err != nil {
		return nil, err
	}

	var maxWorkspaces int
	err = r.db.QueryRowContext(ctx,
		`SELECT max_workspaces FROM user_quotas WHERE user_id = $1`,
		req.UserId,
	).Scan(&maxWorkspaces)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			maxWorkspaces = 3 // default quota
		} else {
			return nil, err
		}
	}

	if count >= maxWorkspaces {
		return nil, ErrQuotaExceeded
	}

	// 2. cek template exists + ambil image
	var templateImage, templateName string
	err = r.db.QueryRowContext(ctx,
		`SELECT name,image FROM templates WHERE id = $1`,
		req.TemplateId,
	).Scan(&templateName, &templateImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}

	// 3. generate namespace: "ws-{userId[:8]}-{name}"
	namespace := GenerateNamespace(req.UserId, req.Name)

	// 4. marshal env_vars
	envJSON, err := json.Marshal(req.EnvVars)
	if err != nil {
		return nil, err
	}

	// 5. insert workspace
	var w Workspace
	var envRaw []byte
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO workspaces (user_id, template_id, name, namespace, status, env_vars)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, template_id, name, namespace, status, env_vars, created_at, updated_at
	`, req.UserId, req.TemplateId, req.Name, namespace, StatusPending, envJSON,
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
	w.Url = fmt.Sprintf("https://%s.wfdnstore.online", w.Id)

	// update url ke db
	_, err = r.db.ExecContext(ctx, `UPDATE workspaces SET url = $1 WHERE id = $2`, w.Url, w.Id)
	if err != nil {
		return nil, err
	}

	PublishEvent(ctx, r.redisClient, WorkspaceJob{
		WorkspaceId: w.Id,
		UserId:      w.UserId,
		Namespace:   w.Namespace,
		Image:       templateImage,
		TemplateId:  req.TemplateId,
		Username:    username,
		Action:      JobCreate,
		EnvVars:     w.EnvVars,
	})

	// 7. invalidate list cache
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
