package workspaceservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (repo *Repository) CreateAddonWorkspace(c context.Context, req CreateWorkspaceAddon) error {
	envVars := make(map[string]any)
	for _, cfg := range req.Config {
		envVars[cfg.Key] = cfg.Value
	}

	data, err := json.Marshal(envVars)
	if err != nil {
		log.Printf("CreateAddonWorkspace: invalid config json: %v", err)
		return fmt.Errorf("invalid addon configuration format")
	}

	var templateAddonId string
	query := `
        INSERT INTO workspace_addons (id, workspace_id, template_addon_id, status, config)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `
	err = repo.db.QueryRowContext(c, query,
		uuid.New().String(),
		req.WorkspaceID,
		req.TemplateAddonId,
		req.Status,
		data,
	).Scan(&templateAddonId)
	if err != nil {
		log.Printf("CreateAddonWorkspace: db error: %v", err)
		return fmt.Errorf("failed to save addon")
	}

	var namespace string
	err = repo.db.QueryRowContext(c, `SELECT namespace FROM workspaces WHERE id = $1`, req.WorkspaceID).Scan(&namespace)
	if err != nil {
		log.Printf("CreateAddonWorkspace: workspace not found: %v", err)
		return fmt.Errorf("workspace not found")
	}

	var image string
	err = repo.db.QueryRowContext(c, "select image from template_addons where id = $1", req.TemplateAddonId).Scan(&image)
	if err != nil {
		log.Printf("CreateAddonWorkspace: workspace not found: %v", err)
		return fmt.Errorf("workspace not found")
	}
	var templateId string
	err = repo.db.QueryRowContext(c, "select template_id from template_addons where id = $1", req.TemplateAddonId).Scan(&templateId)
	if err != nil {
		log.Printf("CreateAddonWorkspace: workspace not found: %v", err)
		return fmt.Errorf("workspace not found")
	}
	debugData, _ := json.Marshal(envVars)
	log.Printf("DEBUG ENV VARS: %s", string(debugData))

	PublishEvent(c, repo.redisClient, WorkspaceJob{
		WorkspaceId: req.WorkspaceID,
		Namespace:   namespace,
		TemplateId:  templateId,
		Action:      JobAdd,
		Image:       image,
		EnvVars:     envVars,
	})

	return nil
}

func (repo *Repository) GetAddonService(c context.Context, workspaceId string) ([]WorkspaceAddon, error) {
	query := `
		SELECT id, workspace_id, template_addon_id, status, config
		FROM workspace_addon
		WHERE workspace_id = $1
	`

	rows, err := repo.db.QueryContext(c, query, workspaceId)
	if err != nil {
		log.Printf("GetAddonService: db error for workspace %s: %v", workspaceId, err)
		return nil, fmt.Errorf("failed to retrieve addons")
	}
	defer rows.Close()

	var addons []WorkspaceAddon
	for rows.Next() {
		var addon WorkspaceAddon
		var configRaw []byte

		if err := rows.Scan(
			&addon.ID,
			&addon.WorkspaceID,
			&addon.TemplateAddonId,
			&addon.Status,
			&configRaw,
		); err != nil {
			log.Printf("GetAddonService: scan error: %v", err)
			return nil, fmt.Errorf("failed to read addon data")
		}

		if err := json.Unmarshal(configRaw, &addon.Config); err != nil {
			log.Printf("GetAddonService: failed to parse config for addon: %v", err)
			return nil, fmt.Errorf("failed to parse addon configuration")
		}

		addons = append(addons, addon)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetAddonService: rows iteration error: %v", err)
		return nil, fmt.Errorf("failed to retrieve addons")
	}

	if len(addons) == 0 {
		return nil, ErrAddonNotFound
	}

	return addons, nil
}
