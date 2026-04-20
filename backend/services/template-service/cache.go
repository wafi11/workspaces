package templateservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wafi11/workspaces/pkg/models"
)

func (r *Repository) setTemplateCache(ctx context.Context, templateId string, t *models.Template) {
	cached := models.CachedTemplate{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
		Category:    t.Category,
		IsPublic:    t.IsPublic,
		Variables:   t.Variables,
		Addons:      t.Addons,
		TemplateUrl: t.TemplateUrl,
		Icon: t.Icon,
		CreatedAt:   t.CreatedAt,
	}
	b, err := json.Marshal(cached)
	if err != nil {
		return
	}
	r.redisClient.Set(ctx, fmt.Sprintf(models.TemplateCacheKey, templateId), b, models.TemplateCacheTTL)
}

func (r *Repository) getTemplateCache(ctx context.Context, templateId string) (*models.Template, error) {
	val, err := r.redisClient.Get(ctx, fmt.Sprintf(models.TemplateCacheKey, templateId)).Bytes()
	if err != nil {
		return nil, err
	}

	var cached models.CachedTemplate
	if err := json.Unmarshal(val, &cached); err != nil {
		return nil, err
	}

	return &models.Template{
		Id:          cached.Id,
		Name:        cached.Name,
		Description: cached.Description,
		Category:    cached.Category,
		IsPublic:    cached.IsPublic,
		Variables:   cached.Variables,
		Addons:      cached.Addons,
		CreatedAt:   cached.CreatedAt,
	}, nil
}