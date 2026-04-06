package userservices

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (r *Repository) GetUserSessions(ctx context.Context, req *GetUserSessionsRequest) (*GetUserSessionsResponse, error) {
	// 1. check cache
	if cached, err := r.getSessionsCache(ctx, req.UserId); err == nil {
		return &GetUserSessionsResponse{Sessions: cached}, nil
	}

	// 2. cache miss → hit db
	query := `
		SELECT id, user_id, is_active, user_agent, ip_address, created_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var (
			s         Session
			createdAt time.Time
		)
		if err := rows.Scan(&s.Id, &s.UserId, &s.IsActive, &s.UserAgent, &s.IpAddress, &createdAt); err != nil {
			return nil, err
		}
		s.CreatedAt = createdAt
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 3. store to cache
	r.setSessionsCache(ctx, req.UserId, sessions)

	return &GetUserSessionsResponse{Sessions: sessions}, nil
}

func (r *Repository) RevokeSession(ctx context.Context, req *RevokeSessionRequest) (*RevokeSessionResponse, error) {
	// fetch user_id dulu untuk invalidate sessions cache
	var userId string
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id FROM sessions WHERE id = $1`, req.SessionId,
	).Scan(&userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE sessions SET is_active = false WHERE id = $1`, req.SessionId,
	)
	if err != nil {
		return nil, err
	}

	// invalidate sessions cache for this user
	r.invalidateSessionsCache(ctx, userId)

	return &RevokeSessionResponse{Message: "session revoked successfully"}, nil
}
