package authservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wafi11/workspaces/config"
)

func (repo *Repository) CreatePAT(c context.Context, req *CreatePATRequest) (*CreatePATResponse, error) {
	raw, hash, err := config.GeneratePAT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate pat: %w", err)
	}

	now := time.Now().UTC()
	fmt.Printf("expires_at : %v", req.ExpiresAt)

	// handle optional expires_at
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := req.ExpiresAt.UTC()
		expiresAt = &t
	}

	query := `
        INSERT INTO personal_access_tokens (
            user_id, name, token_hash, last_used_at, expires_at, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING expires_at
    `

	var dbExpiresAt *time.Time
	err = repo.db.QueryRowContext(c, query,
		req.UserId, req.Name, hash, now, expiresAt, now,
	).Scan(&dbExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pat: %w", err)
	}

	resp := &CreatePATResponse{Token: raw}
	if dbExpiresAt != nil {
		resp.ExpiresAt = dbExpiresAt.UTC().Format(time.RFC3339)
	}

	return resp, nil
}

func (repo *Repository) GetAllPAT(c context.Context, userID string) ([]Pat, error) {
	query := `
		SELECT
			id,
			name,
			last_used_at,
			expires_at,
			created_at
		FROM personal_access_tokens 
		WHERE user_id = $1
	`

	rows, err := repo.db.QueryContext(c, query, userID)
	if err != nil {
		log.Printf("[get pat]  failed to process pat : %s", err.Error())
		return nil, fmt.Errorf("failed to get personal acsess tokens")
	}

	defer rows.Close()

	var data []Pat
	for rows.Next() {
		var pat Pat
		err = rows.Scan(&pat.Id, &pat.Name, &pat.LastUsedAt, &pat.ExpiresAt, &pat.CreatedAt)

		if err != nil {
			log.Printf("[get pat]  failed to process pat : %s", err.Error())

			return nil, fmt.Errorf("failed to find pat : %s", err.Error())
		}

		data = append(data, pat)
	}

	return data, nil
}

func (repo *Repository) DeletePAT(c context.Context, PatId, user_id string) error {
	query := `
		delete from personal_access_tokens where id = $1 and user_id = $2
	`

	_, err := repo.db.ExecContext(c, query, PatId, user_id)
	if err != nil {
		log.Printf("failed to delele pat %s", err.Error())
		return fmt.Errorf("failed to delete pat")
	}

	return nil
}
