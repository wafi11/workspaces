package workspaceservice

import (
	"context"
	"fmt"
	"log"
)


func (repo *Repository) CreateWorkspaceResources(c context.Context, req *CreateWorkspacesResources) error {
	query := `
		INSERT INTO workspace_resources (workspace_id,kind,name,status) values ($1,$2,$3,$4)
	`

	_, err := repo.db.ExecContext(c, query, req.WsID, req.Kind, req.Name, "pending")

	if err != nil {
		log.Printf("failed to create services : %s", err.Error())
		return fmt.Errorf("failed to create services")
	}

	return nil
}

func (repo *Repository) GetWorkspacesResources(c context.Context, WsId string) ([]WorkspacesResources, error) {
	query := `	
		SELECT id,workspace_id,kind,name,status,created_at
		FROM workspace_resources 
		WHERE workspace_id = $1
	`

	rows, err := repo.db.QueryContext(c, query, WsId)

	if err != nil {
		return nil, fmt.Errorf("workspaces not found")
	}

	defer rows.Close()

	var resources []WorkspacesResources
	for rows.Next() {
		var resource WorkspacesResources

		if err := rows.Scan(&resource.Id, &resource.WsID, &resource.Kind, &resource.Name, &resource.Name, &resource.CreatedAt); err != nil {
			log.Printf("failed to scan workspaces resources : %s", err.Error())
			return nil, fmt.Errorf("internal server error")
		}

		resources = append(resources, resource)

	}

	return resources, nil
}
