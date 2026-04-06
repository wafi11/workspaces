package templateservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)


func (r *Repository) CreateTemplateVariable(ctx context.Context, req *CreateVariableRequest,templateId string) (error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO template_variables (id, template_id, key, default_value, required, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New().String(), templateId, req.Key, req.DefaultValue, req.Required, req.Description)

	if err != nil {
		return fmt.Errorf("failed to create template variable: %w", err)
	}

	return nil
}


func (r *Repository) UpdateTemplateVariable(ctx context.Context, id string, req *CreateVariableRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE template_variables
		SET key = $1, default_value = $2, required = $3, description = $4
		WHERE id = $5
	`, req.Key, req.DefaultValue, req.Required, req.Description, id)
	if err != nil {
		return fmt.Errorf("failed to update template variable: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTemplateVariable(ctx context.Context, id string) error{	
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM template_variables
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete template variable: %w", err)
	}
	return nil
}



func (r *Repository) GetTemplateVariables(ctx context.Context, templateID string) ([]TemplateVariable, error) {
	var variables []TemplateVariable
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, template_id, key, default_value, required, description
		FROM template_variables
		WHERE template_id = $1
	`, templateID)
	if err != nil {
		return nil, fmt.Errorf("templates variable not found")
	}
	defer rows.Close()
	for rows.Next() {
		var v TemplateVariable
		if err := rows.Scan(&v.Id, &v.TemplateId, &v.Key, &v.DefaultValue, &v.Required, &v.Description); err != nil {
			return nil, err
		}
		variables = append(variables, v)
	}
	return variables, nil
}