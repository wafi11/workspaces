package templateservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/wafi11/workspaces/pkg/models"
)

func (r *Repository) GetTemplateFiles(ctx context.Context, templateID string) ([]models.TemplateFiles, error) {
	var variables []models.TemplateFiles

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, template_id, filename, sort_order
		FROM template_files
		WHERE template_id = $1
	`, templateID)
	if err != nil {
		return nil, fmt.Errorf("template files not found")
	}
	
	defer rows.Close()
	
	for rows.Next() {
		var v models.TemplateFiles
		if err := rows.Scan(&v.Id, &v.TemplateId, &v.Filename, &v.SortOrder); err != nil {
			return nil, err
		}
		variables = append(variables, v)
	}

	return variables, nil
}

func (r *Repository) CreateTemplateFiles(ctx context.Context, req *models.CreateTemplateFilesRequest, templateId string,tx *sql.Tx) error {
	if tx != nil {
		_, err := tx.ExecContext(ctx, `
		INSERT INTO template_files (id, template_id, filename, sort_order)
		VALUES ($1, $2, $3, $4)
		`, uuid.New().String(), templateId, req.Filename, req.SortOrder)
		if err != nil {
			return fmt.Errorf("failed to create template files: %w", err)
		}
	}else {
		_, err := r.db.ExecContext(ctx, `
		INSERT INTO template_files (id, template_id, filename, sort_order)
		VALUES ($1, $2, $3, $4)
		`, uuid.New().String(), templateId, req.Filename, req.SortOrder)
		if err != nil {
			return fmt.Errorf("failed to create template files: %w", err)
		}
	}
	return nil
}

func (r *Repository) UpdateTemplateFiles(ctx context.Context, id string, req *models.CreateTemplateFilesRequest) error{	
	_, err := r.db.ExecContext(ctx, `
		UPDATE template_files
		SET filename = $1, sort_order = $2
		WHERE id = $3
	`, req.Filename, req.SortOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update template files: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTemplateFiles(ctx context.Context, id string) error{	
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM template_files
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete template files: %w", err)
	}
	return nil
}