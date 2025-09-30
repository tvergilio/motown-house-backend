package db

import (
	"fmt"

	"example.com/web-service-gin/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConnectPostgres returns a *sqlx.DB connected to Postgres using configuration from cfg.
func ConnectPostgres(cfg *config.Config) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("postgres URL is required in config.PostgresURL (set POSTGRES_URL environment variable)")
	}

	db, err := sqlx.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
