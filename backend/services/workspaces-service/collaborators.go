package workspaceservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wafi11/workspaces/config"
	notificationservices "github.com/wafi11/workspaces/services/notifications-service"
)

func (repo *Repository) AddCollaborators(c context.Context, req WorkspaceCollaborator) (*WorkspaceCollaboratorResponse, error) {
	var wsId, userID string
	wsCollId := uuid.New()
	// begin transaction
	tx, err := repo.db.BeginTx(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		select id from users where email = $1
	`

	err = tx.QueryRowContext(c, query, req.Email).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	token, err := config.GenerateTokenWorkspaces(c, &config.TokenWorkspaceRequest{
		UserID:     userID,
		Exp:        timeExpCollaborators,
		AcessLevel: req.Role,
	}, repo.conf)
	if err != nil {
		log.Printf("[add_collaborators] failed to generate token: %s", err.Error())
		return nil, fmt.Errorf("failed to generate invite token")
	}

	queryWs := `SELECT id FROM workspaces WHERE id = $1 and user_id = $2`
	err = repo.db.QueryRowContext(c, queryWs, req.WorkspaceId, req.InvitedBy).Scan(&wsId)
	if err != nil {
		log.Printf("[add_collaborators] failed to find workspace: %s", err.Error())
		return nil, fmt.Errorf("workspace not found")
	}

	queryCreateWsColl := `
		INSERT INTO workspace_collaborators (
			id, workspace_id, user_id, role,
			status, invite_token, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	_, err = tx.ExecContext(c, queryCreateWsColl,
		wsCollId,
		req.WorkspaceId,
		userID,
		req.Role,
		CollaboratorPending,
		token,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Printf("[add_collaborators] failed to create workspace collaborators: %s", err.Error())
		return nil, fmt.Errorf("failed to create workspace collaborators: %w", err)
	}

	metadata, err := json.Marshal(map[string]string{
		"workspace_id": req.WorkspaceId,
		"role":         req.Role,
		"invite_token": token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = repo.notifRepo.CreateNotifications(c, &notificationservices.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationservices.InvitationCollaborator,
		Title:            "Invite Collaboration",
		RetreivedID:      req.InvitedBy,
		Message:          "You are invited to join collaborations",
		Metadata:         metadata,
		IsRead:           false,
	}, tx)
	if err != nil {
		log.Printf("[add_collaborators] failed to create notification: %s", err.Error())
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	repo.hub.SendToUser(userID, map[string]string{
		"message": "You are invited to join collaborations",
		"type":    "notification.unread",
	})

	return &WorkspaceCollaboratorResponse{
		WorkspaceId: wsId,
		Status:      string(CollaboratorPending),
		Token:       token,
	}, nil
}

func (repo *Repository) UpdateCollaborator(c context.Context, req UpdateCollaboratorRequest) error {
	// pastikan yang request adalah owner workspace
	var wsId string
	queryWs := `SELECT id FROM workspaces WHERE id = $1 AND user_id = $2`
	err := repo.db.QueryRowContext(c, queryWs, req.WorkspaceID, req.RequestedBy).Scan(&wsId)
	if err != nil {
		log.Printf("[update_collaborator] workspace not found or unauthorized: %s", err.Error())
		return fmt.Errorf("workspace not found or unauthorized")
	}

	query := `
		UPDATE workspace_collaborators
		SET
			role       = COALESCE(NULLIF($1, ''), role),
			status     = COALESCE(NULLIF($2, ''), status),
			updated_at = $3
		WHERE id = $4 AND workspace_id = $5
	`
	result, err := repo.db.ExecContext(c, query,
		req.Role,
		req.Status,
		time.Now(),
		req.CollaboratorID,
		req.WorkspaceID,
	)
	if err != nil {
		log.Printf("[update_collaborator] failed to update: %s", err.Error())
		return fmt.Errorf("failed to update collaborator")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("collaborator not found")
	}

	return nil
}

func (repo *Repository) RemoveCollaborator(c context.Context, req RemoveCollaboratorRequest) error {
	// boleh remove kalau: dia owner workspace ATAU dia collaborator itu sendiri
	var wsOwnerId string
	queryWs := `SELECT user_id FROM workspaces WHERE id = $1`
	err := repo.db.QueryRowContext(c, queryWs, req.WorkspaceID).Scan(&wsOwnerId)
	if err != nil {
		log.Printf("[remove_collaborator] workspace not found: %s", err.Error())
		return fmt.Errorf("workspace not found")
	}

	var collaboratorUserId string
	queryGetColl := `SELECT user_id FROM workspace_collaborators WHERE id = $1 AND workspace_id = $2`
	err = repo.db.QueryRowContext(c, queryGetColl, req.CollaboratorID, req.WorkspaceID).Scan(&collaboratorUserId)
	if err != nil {
		log.Printf("[remove_collaborator] collaborator not found: %s", err.Error())
		return fmt.Errorf("collaborator not found")
	}

	// authorization check
	isOwner := wsOwnerId == req.RequestedBy
	isSelf := collaboratorUserId == req.RequestedBy
	if !isOwner && !isSelf {
		return fmt.Errorf("unauthorized to remove this collaborator")
	}

	query := `DELETE FROM workspace_collaborators WHERE id = $1 AND workspace_id = $2`
	result, err := repo.db.ExecContext(c, query, req.CollaboratorID, req.WorkspaceID)
	if err != nil {
		log.Printf("[remove_collaborator] failed to delete: %s", err.Error())
		return fmt.Errorf("failed to remove collaborator")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("collaborator not found")
	}

	return nil
}

func (repo *Repository) AcceptOrDeniedInvitationCollborator(ctx context.Context, types, notificationID, userId string) error {
	// 1. Ambil invite_token + collaborator id dari notifikasi
	var inviteToken, userID string
	err := repo.db.QueryRowContext(ctx, `
        SELECT metadata->>'invite_token', user_id 
        FROM notifications WHERE id = $1
    `, notificationID).Scan(&inviteToken, &userID)
	if err != nil {
		log.Printf("[acccept or denied] errr %s", err.Error())
		return fmt.Errorf("notification not found")
	}

	// 2. Validate token
	token, err := config.ValidateTokenWorkspace(inviteToken, repo.conf)
	if err != nil {
		return fmt.Errorf("invite token expired or invalid")
	}

	if token.UserID != userID {
		return fmt.Errorf("failed to update token")
	}

	log.Printf("[accept or denied] inviteToken=%s userID=%s", inviteToken, userID)

	// 3. Update status workspace_collaborators langsung
	newStatus := "denied"
	if types == "accept" {
		newStatus = "active"
	}

	_, err = repo.db.ExecContext(ctx, `
        UPDATE workspace_collaborators 
        SET status = $1, updated_at = NOW()
        WHERE invite_token = $2 AND user_id = $3
    `, newStatus, inviteToken, userID)
	if err != nil {
		return fmt.Errorf("failed to update collaborator status")
	}

	// 4. Mark notif as read
	_, err = repo.db.ExecContext(ctx, `
        UPDATE notifications SET is_read = true WHERE id = $1
    `, notificationID)

	return err
}

func (repo *Repository) GetCollaboratedWorkspaces(ctx context.Context, userID string) ([]CollaboratedWorkspace, error) {
	query := `
        SELECT 
            w.id,
            w.name,
			w.url,
            wc.role,
            wc.status,
            wc.created_at,
			t.name,
			t.icon
        FROM workspace_collaborators wc
        JOIN workspaces w ON w.id = wc.workspace_id
        LEFT JOIN templates t on t.id = w.template_id
		WHERE wc.user_id = $1
          AND wc.status  = 'active'
        ORDER BY wc.created_at DESC
    `

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("failed to get collaborations : %s", err.Error())
		return nil, fmt.Errorf("failed to get collaborated workspaces")
	}
	defer rows.Close()

	var result []CollaboratedWorkspace
	for rows.Next() {
		var cw CollaboratedWorkspace
		err := rows.Scan(&cw.WorkspaceID, &cw.WorkspaceName, &cw.WorkspaceUrl, &cw.Role, &cw.Status, &cw.InvitedAt, &cw.TemplateName, &cw.TemplateIcon)
		if err != nil {
			log.Printf("failed to get collaborations : %s", err.Error())

			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, cw)
	}

	return result, nil
}
