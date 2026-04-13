package workspaceservice

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
	r.redis.Redis.Set(ctx, fmt.Sprintf(workspaceCacheKey, w.Id), b, cacheTTL)
}

func (r *Repository) getWorkspaceCache(ctx context.Context, workspaceId string) (*Workspace, error) {
	val, err := r.redis.Redis.Get(ctx, fmt.Sprintf(workspaceCacheKey, workspaceId)).Bytes()
	if err != nil {
		return nil, err
	}
	var w Workspace
	if err := json.Unmarshal(val, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repository) invalidateWorkspaceCache(ctx context.Context, workspaceId, userId string) {
	r.redis.Redis.Del(ctx, fmt.Sprintf(workspaceCacheKey, workspaceId))
	r.redis.Redis.Del(ctx, fmt.Sprintf(workspacesCacheKey, userId))
}

func GenerateAddonConnectionUrl(addonType AddonUrl, namespace,name,user_id,db_user,db_password,db_name string) string {
    switch addonType {
    case PostgresqlURL :
        return fmt.Sprintf("postgresql://%s:%s@postgres.%s.svc.cluster.local:5432/%s",
            db_user, db_password, namespace, db_name)
    case MysqlURL :
        return fmt.Sprintf("mysql://%s:%s@mysql.%s.svc.cluster.local:3306/%s",
            db_user, db_password, namespace, db_name)
    case RedisURL :
        return fmt.Sprintf("redis://:password@redis.%s.svc.cluster.local:6379",
            namespace)
    default:
        return fmt.Sprintf("%s-%s.wfdnstore.online", name, user_id)
    }
}