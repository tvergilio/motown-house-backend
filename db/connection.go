package db

import (
	"fmt"

	"example.com/web-service-gin/config"
	"github.com/gocql/gocql"
	"github.com/jmoiron/sqlx"
)

// DatabaseConnection holds either a Postgres or Cassandra connection
type DatabaseConnection struct {
	PostgresDB  *sqlx.DB
	CassandraDB *gocql.Session
	Backend     string
}

// Close closes the appropriate database connection
func (dc *DatabaseConnection) Close() error {
	if dc.PostgresDB != nil {
		return dc.PostgresDB.Close()
	}
	if dc.CassandraDB != nil {
		dc.CassandraDB.Close()
		return nil
	}
	return nil
}

// Connect returns a DatabaseConnection based on the configured backend
func Connect(cfg *config.Config) (*DatabaseConnection, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	switch cfg.DBBackend {
	case "postgres":
		db, err := ConnectPostgres(cfg)
		if err != nil {
			return nil, err
		}
		return &DatabaseConnection{
			PostgresDB: db,
			Backend:    "postgres",
		}, nil

	case "cassandra":
		session, err := ConnectCassandra(cfg)
		if err != nil {
			return nil, err
		}
		return &DatabaseConnection{
			CassandraDB: session,
			Backend:     "cassandra",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported database backend: %s", cfg.DBBackend)
	}
}

// ConnectFromEnv loads configuration from environment and returns a connected DB.
func ConnectFromEnv() (*DatabaseConnection, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, err
	}
	return Connect(cfg)
}
