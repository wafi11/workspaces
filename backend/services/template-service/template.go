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

func (r *Repository) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	query := `
		SELECT id, name, description,icon, image, category, is_public, created_at,template_url
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

	var templates []Template
	for rows.Next() {
		var t Template
		if err := rows.Scan(
			&t.Id, &t.Name, &t.Description, &t.Icon,
			&t.Image, &t.Category, &t.IsPublic, &t.CreatedAt, &t.TemplateUrl,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return &ListTemplatesResponse{Templates: templates}, nil
}

func (r *Repository) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {
	// 2. cache miss → hit db
	var t Template
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, image, category, is_public, created_at,template_url,icon
		FROM templates
		WHERE id = $1
	`, req.TemplateId).Scan(
		&t.Id, &t.Name, &t.Description,
		&t.Image, &t.Category, &t.IsPublic, &t.CreatedAt, &t.TemplateUrl, &t.Icon,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}

	return &GetTemplateResponse{Template: &t}, nil
}

func (r *Repository) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. insert template
	var t Template
	err = tx.QueryRowContext(ctx, `
		INSERT INTO templates (name, description, image, category, is_public,icon)
		VALUES ($1, $2, $3, $4, $5,$6)
		RETURNING id, name, description, image, category, is_public, created_at,icon
	`, req.Name, req.Description, req.Image, req.Category, req.IsPublic, req.Icon,
	).Scan(&t.Id, &t.Name, &t.Description, &t.Image, &t.Category, &t.IsPublic, &t.CreatedAt, &req.Icon)

	if err != nil {
		return nil, fmt.Errorf("failed to insert template: %w", err)
	}

	for _, v := range req.Variables {
		r.CreateTemplateVariable(ctx, &v, t.Id)
	}

	for _, v := range req.Files {
		r.CreateTemplateFiles(ctx, &v, t.Id)
	}

	// 3. insert addons
	for _, a := range req.Addons {
		configJSON, err := json.Marshal(a.DefaultConfig)
		if err != nil {
			return nil, err
		}
		a.DefaultConfig = map[string]any{}
		if err := json.Unmarshal(configJSON, &a.DefaultConfig); err != nil {
			return nil, err
		}
		r.CreateTemplateAddon(ctx, &a, t.Id)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// invalidate list cache
	r.redisClient.Del(ctx, templatesCacheKey)

	return &CreateTemplateResponse{
		Template: &t,
		Message:  "template created successfully",
	}, nil
}

func (repo *Repository) UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) error {
	// 1. fetch existing
	existing, err := repo.GetTemplate(ctx, &GetTemplateRequest{TemplateId: id})
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// 2. merge partial request ke existing
	req.merge(existing.Template)

	// 3. update dengan data lengkap — tidak ada mismatch
	query := `
		UPDATE templates 
		SET 
			name         = $1,
			description  = $2,
			image        = $3,
			category     = $4,
			is_public    = $5,
			template_url = $6,
			icon         = $7
		WHERE id = $8
	`
	_, err = repo.db.ExecContext(ctx, query,
		existing.Template.Name,
		existing.Template.Description,
		existing.Template.Image,
		existing.Template.Category,
		existing.Template.IsPublic,
		existing.Template.TemplateUrl,
		existing.Template.Icon,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	repo.redisClient.Del(ctx, templatesCacheKey)

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
	repo.redisClient.Del(ctx, templatesCacheKey)

	return nil
}

func (repo *Repository) GetTemplatesConfigFile(ctx context.Context, id string) ([]TemplateFileInfo, error) {
	var files []TemplateFileInfo

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

func (repo *Repository) GetDetailsInfo(c context.Context, templateId string) (*DetailsInfo, error) {
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

	var detail DetailsInfo
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


func (repo *Repository)  FindTemplateWorkspaceForm(c context.Context)([]TemplateWorkspaceForm,error){
	query := `
		SELECT id,name,icon
		FROM templates
	`

	rows,err := repo.db.QueryContext(c,query)
	if err != nil {
		log.Printf("failed to get template workspaces form")
		return nil,fmt.Errorf("templates not found")
	}

	defer rows.Close()
	
	var templates []TemplateWorkspaceForm
	for rows.Next() {
		var template TemplateWorkspaceForm
		err := rows.Scan(
			&template.ID,template.Name,&template.Icon,
		)

		if err != nil {
			log.Printf("failed to scan template workspace form : %s",err.Error())
			return nil,fmt.Errorf("failed to find templates")
		}

		templates = append(templates,template)
	}

	return templates,nil
}
