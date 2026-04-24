package userservices

import (
	"context"
	"fmt"
	"log"
)

type ProviderUsers struct {
	Name       string `json:"name"`
	ProviderId string `json:"provider_id"`
}

func (repo *Repository) GetAllProvidersUsers(c context.Context, userId string) ([]ProviderUsers, error) {
	query := `
		SELECT name,provider_id FROM providers where user_id = $1
	`

	rows, err := repo.db.DB.QueryContext(c, query, userId)

	if err != nil {
		log.Printf("failed to find providers : %s", err.Error())
		return nil, fmt.Errorf("failed to find providers")
	}

	defer rows.Close()

	var Providers []ProviderUsers
	for rows.Next() {
		var provider ProviderUsers
		err := rows.Scan(&provider.Name, &provider.ProviderId)

		if err != nil {
			return nil, fmt.Errorf("failed to find providers : %s", err.Error())
		}

		Providers = append(Providers, provider)
	}

	return Providers, nil
}
