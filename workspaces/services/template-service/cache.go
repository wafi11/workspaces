package templateservice

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Repository) setTemplateCache(ctx context.Context, templateId string, t *Template) {
	cached := CachedTemplate{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
		Image:       t.Image,
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
	r.redisClient.Set(ctx, fmt.Sprintf(templateCacheKey, templateId), b, cacheTTL)
}

func (r *Repository) getTemplateCache(ctx context.Context, templateId string) (*Template, error) {
	val, err := r.redisClient.Get(ctx, fmt.Sprintf(templateCacheKey, templateId)).Bytes()
	if err != nil {
		return nil, err
	}

	var cached CachedTemplate
	if err := json.Unmarshal(val, &cached); err != nil {
		return nil, err
	}

	return &Template{
		Id:          cached.Id,
		Name:        cached.Name,
		Description: cached.Description,
		Image:       cached.Image,
		Category:    cached.Category,
		IsPublic:    cached.IsPublic,
		Variables:   cached.Variables,
		Addons:      cached.Addons,
		CreatedAt:   cached.CreatedAt,
	}, nil
}