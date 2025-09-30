//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Note: Cassandra integration tests use optimizations (disabled gossip/vnodes) for faster startup.
// Each test starts a fresh Cassandra container which takes ~15-30 seconds to initialize.

// setupTestCassandra spins up a temporary Cassandra container, runs migrations, and returns a connected *gocql.Session and a teardown function.
func setupTestCassandra(t *testing.T) (*gocql.Session, func()) {
	t.Helper()
	ctx := context.Background()

	// Start a new Cassandra container for testing with optimizations for faster startup
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "cassandra:5-jammy",
			ExposedPorts: []string{"9042/tcp"},
			Env: map[string]string{
				// Disable gossip and vnodes for faster startup in testing
				"JVM_EXTRA_OPTS":            "-Dcassandra.skip_wait_for_gossip_to_settle=0 -Dcassandra.load_ring_state=false",
				"CASSANDRA_ENDPOINT_SNITCH": "GossipingPropertyFileSnitch",
				"CASSANDRA_DC":              "datacenter1",
				"CASSANDRA_RACK":            "rack1",
			},
			WaitingFor: wait.ForListeningPort("9042/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Get the host and port to connect to the container
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "9042")
	require.NoError(t, err)

	// Connect to Cassandra
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", host, port.Port()))
	cluster.Consistency = gocql.One // Use ONE instead of QUORUM for single-node testing
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.NumConns = 1                    // Single connection for testing
	cluster.DisableInitialHostLookup = true // Skip host discovery for single node

	// Wait for Cassandra to be ready and create keyspace with exponential backoff
	var session *gocql.Session
	maxRetries := 15
	for i := range maxRetries {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		// Exponential backoff: 1s, 2s, 4s, 8s, then 10s max
		backoff := min(time.Duration(1<<uint(i))*time.Second, 10*time.Second)
		t.Logf("Cassandra connection attempt %d/%d failed: %v. Retrying in %v", i+1, maxRetries, err, backoff)
		time.Sleep(backoff)
	}
	require.NoError(t, err, "Failed to connect to Cassandra after %d retries", maxRetries)

	// Create keyspace with retry logic
	createKeyspaceQuery := "CREATE KEYSPACE IF NOT EXISTS motown WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}"
	for i := range 5 {
		err = session.Query(createKeyspaceQuery).Exec()
		if err == nil {
			break
		}
		t.Logf("Keyspace creation attempt %d failed: %v. Retrying...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, err, "Failed to create keyspace after retries")
	session.Close()

	// Reconnect to the motown keyspace with retry
	cluster.Keyspace = "motown"
	for i := range 5 {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		t.Logf("Keyspace connection attempt %d failed: %v. Retrying...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, err, "Failed to reconnect to motown keyspace")

	// Run migrations by executing the migration files
	err = runCassandraMigrations(session)
	require.NoError(t, err)

	// Return the session and a teardown function to clean up resources
	teardown := func() {
		session.Close()
		_ = container.Terminate(ctx)
	}
	return session, teardown
}

// runCassandraMigrations runs the actual Cassandra migration files
func runCassandraMigrations(session *gocql.Session) error {
	// Find migration file relative to current file's directory
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile)) // Go up two levels from repository/
	migrationFile := filepath.Join(projectRoot, "migrations", "cassandra", "1_create_albums_table.up.cql")

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
	}

	// Split by semicolons and execute each statement
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		if err := session.Query(stmt).Exec(); err != nil {
			return fmt.Errorf("failed to execute statement %q: %w", stmt, err)
		}
	}
	return nil
}

// TestCassandraAlbumRepository_Create tests only the Create method.
func TestCassandraAlbumRepository_Create(t *testing.T) {
	session, teardown := setupTestCassandra(t)
	defer teardown()
	repo := NewCassandraAlbumRepository(session)

	album := Album{Title: "Where Did Our Love Go", Artist: "The Supremes", Price: 9.99, Year: 1964, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music123/v4/5d/c2/4d/5dc24de8-15d7-16e0-7585-72a2bcc721de/14UMGIM62198.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"}
	err := repo.Create(album)
	require.NoError(t, err)

	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, albums, 1)
	require.Equal(t, "Where Did Our Love Go", albums[0].Title)
	require.Equal(t, "The Supremes", albums[0].Artist)
	require.Equal(t, 9.99, albums[0].Price)
	require.Equal(t, 1964, albums[0].Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music123/v4/5d/c2/4d/5dc24de8-15d7-16e0-7585-72a2bcc721de/14UMGIM62198.rgb.jpg/100x100bb.jpg", albums[0].ImageUrl)
	require.Equal(t, "R&B/Soul", albums[0].Genre)
	require.NotEmpty(t, albums[0].ID)
}

// TestCassandraAlbumRepository_GetAll tests only the GetAll method.
func TestCassandraAlbumRepository_GetAll(t *testing.T) {
	session, teardown := setupTestCassandra(t)
	defer teardown()
	repo := NewCassandraAlbumRepository(session)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	_ = repo.Create(Album{Title: "Diana", Artist: "Diana Ross", Price: 2.0, Year: 1980, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/aa/87/1c/aa871c20-95be-38bd-97e3-ecfeb8ec404b/15UMGIM06551.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	_ = repo.Create(Album{Title: "Sex Machine", Artist: "James Brown", Price: 3.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/17/8b/05/178b05de-5855-0136-9827-a0e8a6ccf3db/00602547021656.rgb.jpg/100x100bb.jpg", Genre: "Soul"})

	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, albums, 3)

	// Since Cassandra doesn't guarantee order, we'll check that all albums are present
	titles := make(map[string]bool)
	for _, album := range albums {
		titles[album.Title] = true
		require.NotEmpty(t, album.ID)
		require.NotEmpty(t, album.Artist)
		require.Greater(t, album.Price, 0.0)
		require.Greater(t, album.Year, 0)
		require.NotEmpty(t, album.ImageUrl)
		require.NotEmpty(t, album.Genre)
	}

	require.True(t, titles["ABC"])
	require.True(t, titles["Diana"])
	require.True(t, titles["Sex Machine"])
}

// TestCassandraAlbumRepository_GetByID tests only the GetByID method.
func TestCassandraAlbumRepository_GetByID(t *testing.T) {
	session, teardown := setupTestCassandra(t)
	defer teardown()
	repo := NewCassandraAlbumRepository(session)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, albums)
	id := albums[0].ID

	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "ABC", got.Title)
	require.Equal(t, "Jackson 5", got.Artist)
	require.Equal(t, 1.0, got.Price)
	require.Equal(t, 1970, got.Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", got.ImageUrl)
	require.Equal(t, "R&B/Soul", got.Genre)
	require.Equal(t, id, got.ID)
}

// TestCassandraAlbumRepository_Update tests only the Update method.
func TestCassandraAlbumRepository_Update(t *testing.T) {
	session, teardown := setupTestCassandra(t)
	defer teardown()
	repo := NewCassandraAlbumRepository(session)

	// Create an album
	album := Album{Title: "ABC", Artist: "Shakira", Price: 1.0, Year: 2024, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/8d/97/f4/8d97f427-2d17-1a51-1714-324483eb5fc1/886443546264.jpg/100x100bb.jpg", Genre: "Pop"}
	err := repo.Create(album)
	require.NoError(t, err)

	// Get the created album's ID
	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, albums, 1)
	id := albums[0].ID

	// Update the album
	updated := Album{ID: id, Title: "ABC", Artist: "Jackson 5", Price: 20.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"}
	err = repo.Update(updated)
	require.NoError(t, err)

	// Fetch and verify the updated album
	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "ABC", got.Title)
	require.Equal(t, "Jackson 5", got.Artist)
	require.Equal(t, 20.0, got.Price)
	require.Equal(t, 1970, got.Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", got.ImageUrl)
	require.Equal(t, "R&B/Soul", got.Genre)
	require.Equal(t, id, got.ID)
}

// TestCassandraAlbumRepository_Delete tests only the Delete method.
func TestCassandraAlbumRepository_Delete(t *testing.T) {
	session, teardown := setupTestCassandra(t)
	defer teardown()
	repo := NewCassandraAlbumRepository(session)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, albums)
	id := albums[0].ID

	err = repo.Delete(id)
	require.NoError(t, err)

	albums, err = repo.GetAll()
	require.NoError(t, err)
	require.Len(t, albums, 0)
}
