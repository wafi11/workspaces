package userservices

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/wafi11/workspaces/pkg/models"
)



func (repo *Repository) GetUserQuota(ctx context.Context,userId string) (*models.UserQuota, error) {
	var quota models.UserQuota
	query := `
		SELECT id, user_id, max_workspaces, max_storage_gb, max_ram_mb, max_cpu_cores
		FROM user_quotas
		WHERE user_id = $1
	`

	err := repo.db.GetContext(ctx,&quota, query, userId)
	if err != nil {
		log.Printf("Error fetching user quota for user_id %s: %v", userId, err)
		return nil, fmt.Errorf("user not found"	)
	}

	return &quota, nil
}

func (repo *Repository) UpdateUserQuota(ctx context.Context,quota *models.UserQuota) error {
	query := `
		UPDATE user_quotas
		SET max_workspaces = $1, max_storage_gb = $2, max_ram_mb = $3, max_cpu_cores = $4
		WHERE user_id = $5
	`
	_, err := repo.db.ExecContext(ctx,query, quota.MaxWorkspaces, quota.MaxStorageGB, quota.MaxRamMB, quota.MaxCpuCores, quota.UserID)
	log.Printf("Updated user quota for user_id %s: %v", quota.UserID, err)
	return fmt.Errorf("failed to update user quota for user_id %s: %w", quota.UserID, err)
}

func (repo *Repository) CreateUserQuota(ctx context.Context,quota *models.UserQuota) error {
	quota.ID = uuid.New().String()
	query := `
		INSERT INTO user_quotas (id, user_id, max_workspaces, max_storage_gb, max_ram_mb, max_cpu_cores)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := repo.db.ExecContext(ctx,query, quota.ID, quota.UserID, quota.MaxWorkspaces, quota.MaxStorageGB, quota.MaxRamMB, quota.MaxCpuCores)
	log.Printf("Created user quota for user_id %s: %v", quota.UserID, err)
	return fmt.Errorf("user not found")
}