package userservices

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Repository) setUserCache(ctx context.Context, userId string, u *User) {
	data := CachedUser{
		Id:          u.Id,
		Username:    u.Username,
		Email:       u.Email,
		TerminalUrl: *u.TerminalUrl,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	r.redisClient.Set(ctx, fmt.Sprintf(userCacheKey, userId), b, cacheTTL)
}

func (r *Repository) getUserCache(ctx context.Context, userId string) (*User, error) {
	val, err := r.redisClient.Get(ctx, fmt.Sprintf(userCacheKey, userId)).Bytes()
	if err != nil {
		return nil, err
	}

	var cached CachedUser
	if err := json.Unmarshal(val, &cached); err != nil {
		return nil, err
	}

	return &User{
		Id:          cached.Id,
		Username:    cached.Username,
		Email:       cached.Email,
		TerminalUrl: &cached.TerminalUrl,
		CreatedAt:   cached.CreatedAt,
		UpdatedAt:   cached.UpdatedAt,
	}, nil
}

func (r *Repository) invalidateUserCache(ctx context.Context, userId string) {
	r.redisClient.Del(ctx, fmt.Sprintf(userCacheKey, userId))
}

func (r *Repository) setSessionsCache(ctx context.Context, userId string, sessions []Session) {
	var data []CachedSession
	for _, s := range sessions {
		data = append(data, CachedSession{
			Id:        s.Id,
			UserId:    s.UserId,
			IsActive:  s.IsActive,
			UserAgent: s.UserAgent,
			IpAddress: s.IpAddress,
			CreatedAt: s.CreatedAt,
		})
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	r.redisClient.Set(ctx, fmt.Sprintf(sessionCacheKey, userId), b, cacheTTL)
}

func (r *Repository) getSessionsCache(ctx context.Context, userId string) ([]Session, error) {
	val, err := r.redisClient.Get(ctx, fmt.Sprintf(sessionCacheKey, userId)).Bytes()
	if err != nil {
		return nil, err // cache miss
	}

	var data []CachedSession
	if err := json.Unmarshal(val, &data); err != nil {
		return nil, err
	}

	var sessions []Session
	for _, s := range data {
		sessions = append(sessions, Session{
			Id:        s.Id,
			UserId:    s.UserId,
			IsActive:  s.IsActive,
			UserAgent: s.UserAgent,
			IpAddress: s.IpAddress,
			CreatedAt: s.CreatedAt,
		})
	}
	return sessions, nil
}

func (r *Repository) invalidateSessionsCache(ctx context.Context, userId string) {
	r.redisClient.Del(ctx, fmt.Sprintf(sessionCacheKey, userId))
}
