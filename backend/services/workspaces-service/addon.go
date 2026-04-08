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

	// --- 1. SETUP DEFAULT RESOURCE ---
	// Sementara kita hardcode dulu nilainya
	const (
		defaultAddonCPU = "0.20"
		defaultAddonMem = "128Mi"
		defaultAddonStg = "1Gi"
	)

	// 2. Insert ke workspace_addons
	var addonId string = uuid.New().String()
	query := `
        INSERT INTO workspace_addons (id, workspace_id, template_addon_id, status, config)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = repo.db.ExecContext(c, query,
		addonId,
		req.WorkspaceID,
		req.TemplateAddonId,
		"pending",
		data,
	)
	if err != nil {
		log.Printf("CreateAddonWorkspace: db error: %v", err)
		return fmt.Errorf("failed to save addon")
	}

	// 3. Ambil data Namespace, Image, dan TemplateID (Gabung jadi satu query biar efisien)
	var namespace, image, name, templateId string
	queryDetails := `
        SELECT w.namespace,w.name, ta.image, ta.template_id 
        FROM workspaces w
        JOIN template_addons ta ON ta.id = $2
        WHERE w.id = $1
    `
	err = repo.db.QueryRowContext(c, queryDetails, req.WorkspaceID, req.TemplateAddonId).Scan(
		&namespace, &name, &image, &templateId,
	)
	if err != nil {
		log.Printf("CreateAddonWorkspace: lookup error: %v", err)
		return fmt.Errorf("workspace or addon template not found")
	}

	// 4. Catat ke workspace_resources (Gunakan DEFAULT)
	// Ini penting biar di K8S Operator nanti bisa dapet angka resources-nya
	_, _ = repo.db.ExecContext(c, `
        INSERT INTO workspace_resources (workspace_id, kind, name, cpu_cores, ram_mb, storage_gb, status)
        VALUES ($1, 'addon', $2, $3, $4, $5, 'pending')`,
		req.WorkspaceID, "addon-"+addonId[:8], 0.20, 128, 1,
	)

	log.Printf("DEBUG ADDON: Image=%s, Namespace=%s, CPU=%s", image, namespace, defaultAddonCPU)

	// 5. Publish Event
	PublishEvent(c, repo.redisClient, WorkspaceJob{
		WorkspaceId:    req.WorkspaceID,
		Namespace:      namespace,
		TemplateId:     templateId,
		Action:         JobAdd,
		Image:          image,
		EnvVars:        envVars,
		Name:           name,
		CPURequest:     defaultAddonCPU,
		MemoryRequest:  defaultAddonMem,
		StorageRequest: defaultAddonStg,
		CPULimit:       "1",
		MemoryLimit:    defaultAddonMem,
	})

	return nil
}

func (repo *Repository) GetAddonService(c context.Context, workspaceId string) ([]WorkspaceAddonDetails, error) {
	query := `
		SELECT 
			wa.id, 
			t.name,
			t.icon,
			wa.status, 
			wa.config
		FROM workspace_addons wa
		LEFT JOIN template_addons ta on  ta.id = wa.template_addon_id 
		LEFT JOIN templates t on t.id = ta.template_id
		WHERE wa.workspace_id = $1
	`

	rows, err := repo.db.QueryContext(c, query, workspaceId)
	if err != nil {
		log.Printf("GetAddonService: db error for workspace %s: %v", workspaceId, err)
		return nil, fmt.Errorf("failed to retrieve addons")
	}
	defer rows.Close()

	var addons []WorkspaceAddonDetails
	for rows.Next() {
		var addon WorkspaceAddonDetails

		if err := rows.Scan(
			&addon.ID,
			&addon.Name,
			&addon.Icon,
			&addon.Status,
			&addon.Config,
		); err != nil {
			log.Printf("GetAddonService: scan error: %v", err)
			return nil, fmt.Errorf("failed to read addon data")
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
