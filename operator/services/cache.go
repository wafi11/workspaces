package services

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Repository) setWorkspaceCache(ctx context.Context, w *Workspace) {
	b, err := json.Marshal(w)
	if err != nil {
		return
	}
	r.redisClient.Set(ctx, fmt.Sprintf(workspaceCacheKey, w.Id), b, cacheTTL)
}

func (r *Repository) getWorkspaceCache(ctx context.Context, workspaceId string) (*Workspace, error) {
	val, err := r.redisClient.Get(ctx, fmt.Sprintf(workspaceCacheKey, workspaceId)).Bytes()
	if err != nil {
		return nil, err
	}
	var w Workspace
	if err := json.Unmarshal(val, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

