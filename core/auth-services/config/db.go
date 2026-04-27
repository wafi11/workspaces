package config

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DBConn(conStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("connected to database successfully")
	return db, nil
}
