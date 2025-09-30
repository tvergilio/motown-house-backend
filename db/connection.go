package db

import (
	"fmt"

	"example.com/web-service-gin/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect returns a *sqlx.DB connected to Postgres using configuration from cfg.

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	var dsn string
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("postgres URL is required in config.PostgresURL")
	}
	dsn = cfg.PostgresURL

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// ConnectFromEnv loads configuration from environment and returns a connected DB.
// Convenience wrapper for callers that prefer a no-argument call.
func ConnectFromEnv() (*sqlx.DB, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, err
	}
	return Connect(cfg)
}
