package workspaceservice

import (
	"context"
	"fmt"
	"log"
	"time"
)

type WorkspacePort struct {
	ID          string    `json:"id.omitempty" db:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty" db:"workspace_id"`
	Port        int       `json:"port" db:"port"`
	SubDomain   string    `json:"subdomain" db:"subdomain"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (repo *Repository) ListWorkspacePorts(ctx context.Context, workspaceID string) ([]WorkspacePort, error) {
	query := `
		SELECT id,port, subdomain, created_at
		FROM workspace_ports
		WHERE workspace_id = $1
		ORDER BY created_at DESC
	`
	var ports []WorkspacePort
	if err := repo.db.SelectContext(ctx, &ports, query, workspaceID); err != nil {
		return nil, err
	}
	return ports, nil
}

func (repo *Repository) CreateWorkspacePort(ctx context.Context, workspaceID string, port int, userID string) (*WorkspacePort, error) {
	var url string

	querySelect := `
		select url from workspaces where id = $1 and user_id = $2
	`

	if err := repo.db.QueryRowContext(ctx,querySelect,workspaceID,userID).Scan(&url); err != nil {
		log.Printf("workspace err : %s",err.Error())
		return nil,fmt.Errorf("workspace not found")
	}
	
	subDomain := fmt.Sprintf("%d-%s",port,url)

	query := `
		INSERT INTO workspace_ports (id, workspace_id, port, subdomain, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())
		RETURNING id, workspace_id, port, subdomain, created_at
	`
	var p WorkspacePort
	if err := repo.db.QueryRowContext(ctx, query, workspaceID, port, subDomain).Scan(
		&p.ID, &p.WorkspaceID, &p.Port, &p.SubDomain, &p.CreatedAt,
	); err != nil {
				log.Printf("workspace err : %s",err.Error())

		return nil, err
	}
	return &p, nil
}

func (repo *Repository) DeleteWorkspacePort(ctx context.Context, workspaceID string, port int) error {
	query := `
		DELETE FROM workspace_ports
		WHERE workspace_id = $1 AND port = $2
	`
	_, err := repo.db.ExecContext(ctx, query, workspaceID, port)
	return err
}

func (repo *Repository) GetWorkspacePort(ctx context.Context, workspaceID string, port int) (*WorkspacePort, error) {
	query := `
		SELECT id, workspace_id, port, subdomain, created_at
		FROM workspace_ports
		WHERE workspace_id = $1 AND port = $2
	`
	var p WorkspacePort
	if err := repo.db.QueryRowContext(ctx, query, workspaceID, port).Scan(
		&p.ID, &p.WorkspaceID, &p.Port, &p.SubDomain, &p.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo *Repository) ValidateWorkspaceOwner(ctx context.Context, workspaceID, userID string) error {
	var id string
	err := repo.db.QueryRowContext(ctx,
		`SELECT id FROM workspaces WHERE id = $1 AND user_id = $2`,
		workspaceID, userID,
	).Scan(&id)
	if err != nil {
		return ErrWorkspaceNotFound
	}
	return nil
}