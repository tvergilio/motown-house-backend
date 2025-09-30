package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	// DBBackend chooses which database backend to use: "postgres" or "cassandra".
	DBBackend string

	// Postgres specific (prefer POSTGRES_URL if present)
	// POSTGRES_URL is the preferred single-variable form (keeps the repository's
	// previous behaviour of building a connection string from parts if unset).
	PostgresURL string

	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string

	// Cassandra specific
	CassandraHosts    string // comma-separated
	CassandraKeyspace string
}

// LoadFromEnv reads environment variables and returns a Config.
// It supports the new variables and falls back to legacy POSTGRES_* vars
// so this change remains compatible with existing env files.
func LoadFromEnv() (*Config, error) {
	c := &Config{
		DBBackend:         strings.ToLower(strings.TrimSpace(os.Getenv("DB_BACKEND"))),
		PostgresURL:       strings.TrimSpace(os.Getenv("POSTGRES_URL")),
		CassandraHosts:    strings.TrimSpace(os.Getenv("CASSANDRA_HOSTS")),
		CassandraKeyspace: strings.TrimSpace(os.Getenv("CASSANDRA_KEYSPACE")),
	}

	// sensible defaults
	if c.DBBackend == "" {
		c.DBBackend = "postgres"
	}

	// Validate backend
	if c.DBBackend != "postgres" && c.DBBackend != "cassandra" {
		return nil, errors.New("DB_BACKEND must be either 'postgres' or 'cassandra'")
	}

	// For postgres we require a single connection URL to keep behaviour explicit.
	if c.DBBackend == "postgres" {
		if c.PostgresURL == "" {
			return nil, errors.New("POSTGRES_URL must be set when DB_BACKEND=postgres")
		}
	}

	// Normalise cassandra hosts (ensure comma separated if space separated)
	if c.CassandraHosts != "" {
		c.CassandraHosts = strings.ReplaceAll(c.CassandraHosts, " ", ",")
	}

	return c, nil
}
