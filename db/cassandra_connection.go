package db

import (
	"fmt"
	"strings"
	"time"

	"example.com/web-service-gin/config"
	"github.com/gocql/gocql"
)

// ConnectCassandra returns a *gocql.Session connected to Cassandra using configuration from cfg.
func ConnectCassandra(cfg *config.Config) (*gocql.Session, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Default values for local development
	hosts := "localhost:9042"
	keyspace := "motown"

	if cfg.CassandraHosts != "" {
		hosts = cfg.CassandraHosts
	}
	if cfg.CassandraKeyspace != "" {
		keyspace = cfg.CassandraKeyspace
	}

	// Parse hosts (handle comma-separated values)
	hostList := strings.Split(hosts, ",")
	for i, host := range hostList {
		hostList[i] = strings.TrimSpace(host)
	}

	// Create cluster configuration
	cluster := gocql.NewCluster(hostList...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra cluster %v: %w", hostList, err)
	}

	return session, nil
}
