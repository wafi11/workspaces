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
	"github.com/wafi11/workspaces/pkg/utils"
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

	if req.Password == "" {
		return nil, fmt.Errorf("password must be required")
	}
	hashedPassword, err := utils.HashPassword(req.Password)

	envJSON, _ := json.Marshal(req.EnvVars)
	err = tx.QueryRowContext(ctx, `
        INSERT INTO workspaces (user_id, name, status, env_vars,password)
        VALUES ($1, $2, $3, $4,$5)
        RETURNING id, user_id, name, status`,
		req.UserId, req.Name, StatusPending, envJSON, hashedPassword,
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
		`, w.Id, w.Name, w.Status, req.ReqCpuCores, req.ReqRam, req.LimitRam, req.LimitCpuCores,
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
		publisher.PublishEvent(ctx, r.redisClient, publisher.WorkspaceJob{
			UserId:        w.UserId,
			TemplateId:    req.TemplateId,
			Name:          req.Name,
			Username:      username,
			Action:        publisher.JobAdd,
			EnvVars:       req.EnvVars,
			CPURequest:    fmt.Sprintf("%.2f", req.ReqCpuCores),
			MemoryRequest: fmt.Sprintf("%dMi", req.ReqRam),
			CPULimit:      fmt.Sprintf("%.2f", req.LimitCpuCores),
			MemoryLimit:   fmt.Sprintf("%dMi", req.LimitRam),
		})
	}

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
		SELECT id, user_id, name,status,url, env_vars, created_at, updated_at
		FROM workspaces
		WHERE id = $1
	`, req.WorkspaceId,
	).Scan(
		&w.Id, &w.UserId,
		&w.Name, &w.Status, &w.Url,
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var res struct {
		UserId     string
		LimitRAM   int
		LimitCPU   float64
		ReqRAM     int
		ReqCPU     float64
		CurrStatus WorkspaceStatus
	}
	err = tx.QueryRowContext(ctx, `
		SELECT w.user_id, w.status,
			wr.limit_ram_mb, wr.limit_cpu_cores,
			wr.ram_mb_req, wr.cpu_cores_req
		FROM workspaces w
		JOIN workspace_resources wr ON wr.workspace_id = w.id
		WHERE w.id = $1
		FOR UPDATE`, workspaceId,
	).Scan(&res.UserId, &res.CurrStatus, &res.LimitRAM, &res.LimitCPU, &res.ReqRAM, &res.ReqCPU)
	if err != nil {
		return fmt.Errorf("workspace tidak ditemukan: %w", err)
	}

	// Guard: jangan proses kalau status sama
	if res.CurrStatus == status {
		return fmt.Errorf("workspace sudah dalam status %s", status)
	}

	if status == "running" {
		// Ambil quota user + lock
		var maxCPU, usedCPU float64
		var maxRAM, usedRAM int
		err = tx.QueryRowContext(ctx, `
			SELECT max_cpu_cores, used_cpu_cores, max_ram_mb, used_ram_mb
			FROM user_quotas WHERE user_id = $1 FOR UPDATE`, res.UserId,
		).Scan(&maxCPU, &usedCPU, &maxRAM, &usedRAM)
		if err != nil {
			return fmt.Errorf("gagal ambil quota: %w", err)
		}

		// Validasi resource
		if usedCPU+res.ReqCPU > maxCPU {
			return fmt.Errorf("quota CPU tidak cukup, matikan workspace lain dulu")
		}
		if usedRAM+res.ReqRAM > maxRAM {
			return fmt.Errorf("quota RAM tidak cukup, matikan workspace lain dulu")
		}

		// Claim resource
		_, err = tx.ExecContext(ctx, `
			UPDATE user_quotas
			SET used_cpu_cores = used_cpu_cores + $1,
				used_ram_mb = used_ram_mb + $2
			WHERE user_id = $3`,
			res.ReqCPU, res.ReqRAM, res.UserId,
		)
		if err != nil {
			return fmt.Errorf("gagal update quota: %w", err)
		}
	}

	if status == "stopped" {
		_, err = tx.ExecContext(ctx, `
			UPDATE user_quotas
			SET used_cpu_cores = GREATEST(used_cpu_cores - $1, 0),
				used_ram_mb = GREATEST(used_ram_mb - $2, 0)
			WHERE user_id = $3`,
			res.ReqCPU, res.ReqRAM, res.UserId,
		)
		if err != nil {
			return fmt.Errorf("gagal release quota: %w", err)
		}
	}

	// Update status workspace
	_, err = tx.ExecContext(ctx,
		`UPDATE workspaces SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, workspaceId,
	)
	if err != nil {
		return fmt.Errorf("gagal update status: %w", err)
	}

	return tx.Commit()
}
