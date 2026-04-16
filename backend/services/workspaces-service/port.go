package workspaceservice

import (
	"context"
	"fmt"
	"log"
	"time"

	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
	"github.com/wafi11/workspaces/pkg/proto"
)

type WorkspacePort struct {
	ID          string    `json:"id,omitempty" db:"id"`
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
	var url,workspace_name string

	querySelect := `
		select url,name from workspaces where id = $1 and user_id = $2
	`

	if err := repo.db.QueryRowContext(ctx,querySelect,workspaceID,userID).Scan(&url,&workspace_name); err != nil {
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


	messagebroker.PublishEvent(ctx,repo.redis.Redis,&proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_CreatePort{
				CreatePort: &proto.CreatePort{
					Port: int32(port),
					Domain: subDomain,
					WorkspaceName: workspace_name,
					UserId: userID,
				},
		},
	})

	return &p, nil
}

func (repo *Repository) DeleteWorkspacePort(ctx context.Context, workspaceID string, port int) error {
	var workspace_name,userId string

	queryWs := `
		SELECT
			name,
			user_id
		FROM workspaces 
		where id = $1
	`

	err := repo.db.QueryRowContext(ctx,queryWs,workspaceID).Scan(&workspace_name,&userId)

	if err != nil {
		log.Printf("[workspace port] get details workspace error : %s",err.Error())
		return fmt.Errorf("workspace not found")
	}

	query := `
		DELETE FROM workspace_ports
		WHERE workspace_id = $1 AND port = $2
	`
	_, err = repo.db.ExecContext(ctx, query, workspaceID, port)

	if err != nil {
		log.Printf("[workspace port] DELETE workspace error : %s",err.Error())
		return fmt.Errorf("workspaces not found")
	}

	messagebroker.PublishEvent(ctx,repo.redis.Redis,&proto.WorkspaceEnvelope{
		Payload: &proto.WorkspaceEnvelope_DeletePort{
			DeletePort: &proto.DeletePort{
				Port: int32(port),
				WorkspaceName: workspace_name,
				UserId: userId,
			},
		},
	})

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