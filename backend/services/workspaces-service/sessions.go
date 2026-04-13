package workspaceservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
)

func (repo *Repository) CreateWorkspaceSessions(ctx context.Context, req CreateWorkspaceSessions) error {
	now := time.Now().UTC()
	expiresAt := now.Add(5 * time.Minute)
	

	query := `
        INSERT INTO workspace_sessions (
            workspace_id, user_id, status, started_at, expires_at,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := repo.db.ExecContext(ctx, query, req.WorkspaceId, req.UserId, req.Status, now, expiresAt, now, now)
	if err != nil {
		log.Printf("failed to start sessions: %s", err.Error())
		return fmt.Errorf("failed to start sessions")
	}

	return nil
}

func (repo *Repository) HandleStopWorkspace() asynq.HandlerFunc {
    return func(ctx context.Context, t *asynq.Task) error {
        var req messagebroker.TaskSchedulling
        if err := json.Unmarshal(t.Payload(), &req); err != nil {
            return err
        }

        return repo.AutoStopWorkspace(ctx,req.WorkspaceID)
    }
}



func (repo *Repository) CanStartWorkspace(ctx context.Context, workspaceID string,tx *sql.Tx) (bool, error) {
    query := `
        SELECT expires_at FROM workspace_sessions
        WHERE workspace_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `
    var nextStartAt *time.Time
    err := tx.QueryRowContext(ctx, query, workspaceID).Scan(&nextStartAt)
    if err == sql.ErrNoRows {
        return true, nil
    }
    if err != nil {
        return false, err
    }
    if nextStartAt == nil {
        return true, nil
    }
    return time.Now().UTC().After(*nextStartAt), nil
}


