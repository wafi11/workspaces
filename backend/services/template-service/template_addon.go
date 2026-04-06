package templateservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) CreateTemplateAddon(ctx context.Context, req *CreateAddonRequest, templateId string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO template_addons (id, template_id, name, description, default_config)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New().String(), templateId, req.Name, req.Description, req.DefaultConfig)
	if err != nil {
		return fmt.Errorf("failed to create template addon: %w", err)
	}
	return nil
}

func (r *Repository) GetTemplateAddons(ctx context.Context, templateID string) ([]TemplateAddon, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, default_config
		FROM template_addons
		WHERE template_id = $1
	`, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template addons: %w", err)
	}
	defer rows.Close()

	var addons []TemplateAddon
	for rows.Next() {
		var addon TemplateAddon

		if err := rows.Scan(&addon.Id, &addon.Name, &addon.Description, &addon.DefaultConfig); err != nil {
			return nil, fmt.Errorf("failed to scan template addon: %w", err)
		}
		addons = append(addons, addon)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating template addons: %w", err)
	}
	return addons, nil
}

func (r *Repository) UpdateTemplateAddon(ctx context.Context, id string, req *CreateAddonRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE template_addons
		SET name = $1, description = $2, default_config = $3
		WHERE id = $4
	`, req.Name, req.Description, req.DefaultConfig, id)
	if err != nil {
		return fmt.Errorf("failed to update template addon: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTemplateAddon(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM template_addons
		WHERE id = $1
		`, id)
	if err != nil {
		return fmt.Errorf("failed to delete template addon: %w", err)
	}
	return nil
}