package templateservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/wafi11/workspaces/pkg/models"
)

type Repository struct {
	db          *sqlx.DB
	minioClient *minio.Client
	redisClient *redis.Client
}

func NewRepository(db *sqlx.DB, redis *redis.Client, minioClient *minio.Client) *Repository {
	return &Repository{
		db:          db,
		redisClient: redis,
		minioClient: minioClient,
	}
}

func (r *Repository) ListTemplates(ctx context.Context, req *models.ListTemplatesRequest) (*models.ListTemplatesResponse, error) {
	query := `
		SELECT id, name, description,icon, category, is_public, created_at,template_url
		FROM templates
		WHERE is_public = true
	`
	args := []any{}

	if req.Category != "" {
		query += ` AND category = $1`
		args = append(args, req.Category)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var t models.Template
		if err := rows.Scan(
			&t.Id, &t.Name, &t.Description, &t.Icon,
			&t.Category, &t.IsPublic, &t.CreatedAt, &t.TemplateUrl,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return &models.ListTemplatesResponse{Templates: templates}, nil
}

func (r *Repository) GetTemplate(ctx context.Context, req *models.GetTemplateRequest) (*models.GetTemplateResponse, error) {
	// 2. cache miss → hit db
	var t models.Template
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, category, is_public, created_at,template_url,icon
		FROM templates
		WHERE id = $1
	`, req.TemplateId).Scan(
		&t.Id, &t.Name, &t.Description, &t.Category, &t.IsPublic, &t.CreatedAt, &t.TemplateUrl, &t.Icon,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}

	return &models.GetTemplateResponse{Template: &t}, nil
}

func (r *Repository) CreateTemplate(ctx context.Context, req *models.CreateTemplateRequest) (*models.CreateTemplateResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. insert template
	var t models.Template
	err = tx.QueryRowContext(ctx, `
		INSERT INTO templates (name, description, category, is_public,icon,template_url)
		VALUES ($1, $2, $3, $4, $5,$6)
		RETURNING id, name, description, category, is_public, created_at,icon
	`, req.Name, req.Description, req.Category, req.IsPublic, req.Icon, fmt.Sprintf("templates/%s", req.Name),
	).Scan(&t.Id, &t.Name, &t.Description, &t.Category, &t.IsPublic, &t.CreatedAt, &req.Icon)

	if err != nil {
		log.Printf("error create templates : %s", err)
		return nil, fmt.Errorf("failed to insert template: %w", err)
	}

	for _, v := range req.Variables {
		err = r.CreateTemplateVariable(ctx, &v, t.Id, tx.Tx)
		if err != nil {
			log.Printf("insert template error  : %s", err.Error())
			return nil, fmt.Errorf("failed to create variables template")
		}
	}

	for _, v := range req.Files {
		err = r.CreateTemplateFiles(ctx, &v, t.Id, tx.Tx)

		if err != nil {
			log.Printf("insert template files errr  : %s", err.Error())

			return nil, fmt.Errorf("failed to insert addon  template")
		}
	}

	// 3. insert addons
	for _, a := range req.Addons {

		if err := r.CreateTemplateAddon(ctx, &a, t.Id, tx.Tx); err != nil {
			log.Printf("insert template addon error: %s", err.Error())
			return nil, fmt.Errorf("failed to create addon template")
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// invalidate list cache
	r.redisClient.Del(ctx, models.TemplatesCacheKey)

	return &models.CreateTemplateResponse{
		Template: &t,
		Message:  "template created successfully",
	}, nil
}

func (repo *Repository) UpdateTemplate(ctx context.Context, id string, req *models.UpdateTemplateRequest) error {
	// 1. fetch existing
	existing, err := repo.GetTemplate(ctx, &models.GetTemplateRequest{TemplateId: id})
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// 2. merge partial request ke existing
	req.Merge(existing.Template)

	// 3. update dengan data lengkap — tidak ada mismatch
	query := `
		UPDATE templates 
		SET 
			name         = $1,
			description  = $2,
			category     = $3,
			is_public    = $4,
			template_url = $5,
			icon         = $6
		WHERE id = $7
	`
	_, err = repo.db.ExecContext(ctx, query,
		existing.Template.Name,
		existing.Template.Description,
		existing.Template.Category,
		existing.Template.IsPublic,
		existing.Template.TemplateUrl,
		existing.Template.Icon,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	repo.redisClient.Del(ctx, models.TemplatesCacheKey)

	return nil
}
func (repo *Repository) DeleteTemplate(ctx context.Context, id string) error {
	query := `
		DELETE FROM templates where id = $1
	`
	_, err := repo.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete workspaces")
	}
	repo.redisClient.Del(ctx, models.TemplatesCacheKey)

	return nil
}

func (repo *Repository) GetTemplatesConfigFile(ctx context.Context, id string) ([]models.TemplateFileInfo, error) {
	var files []models.TemplateFileInfo

	query := `
        SELECT 
			t.template_url,
            tf.filename,
            tf.sort_order
        FROM template_files tf
        INNER JOIN templates t ON tf.template_id = t.id
        WHERE t.id = $1
        ORDER BY tf.sort_order ASC
    `

	err := repo.db.SelectContext(ctx, &files, query, id)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (repo *Repository) GetDetailsInfo(c context.Context, templateId string) (*models.DetailsInfo, error) {
	query := `
        SELECT t.name, (
            SELECT json_agg(
                json_build_object('id', ta.id, 'name', ta.name)
                ORDER BY ta.name ASC
            )
            FROM template_addons ta
            WHERE ta.template_id = t.id
        ) as addon_list, (
            SELECT json_agg(
                json_build_object('key', tv."key", 'required', tv.required)
                ORDER BY tv."key" ASC
            )
            FROM template_variables tv
            WHERE tv.template_id = t.id
        ) as variables
        FROM templates t
        WHERE t.id = $1
    `

	var detail models.DetailsInfo
	var addonList json.RawMessage
	var varJSON json.RawMessage

	err := repo.db.QueryRowContext(c, query, templateId).Scan(&detail.TemplateName, &addonList, &varJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		log.Printf("failed to scan template: %s", err)
		return nil, fmt.Errorf("failed to find template")
	}

	if addonList != nil {
		if err := json.Unmarshal(addonList, &detail.Addons); err != nil {
			log.Printf("failed to unmarshal variables: %s", err)
			return nil, fmt.Errorf("failed to find template")
		}
	}
	if varJSON != nil {
		if err := json.Unmarshal(varJSON, &detail.Variables); err != nil {
			log.Printf("failed to unmarshal variables: %s", err)
			return nil, fmt.Errorf("failed to find template")
		}
	}

	return &detail, nil
}

func (repo *Repository) FindTemplateWorkspaceForm(c context.Context) ([]models.TemplateWorkspaceForm, error) {
	query := `
		SELECT id,name,icon
		FROM templates
		WHERE is_public = true
	`

	rows, err := repo.db.QueryContext(c, query)
	if err != nil {
		log.Printf("failed to get template workspaces form")
		return nil, fmt.Errorf("templates not found")
	}

	defer rows.Close()

	var templates []models.TemplateWorkspaceForm
	for rows.Next() {
		var template models.TemplateWorkspaceForm
		err := rows.Scan(
			&template.ID, &template.Name, &template.Icon,
		)

		if err != nil {
			log.Printf("failed to scan template workspace form : %s", err.Error())
			return nil, fmt.Errorf("failed to find templates")
		}

		templates = append(templates, template)
	}

	return templates, nil
}
