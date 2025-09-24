package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// migrationsPath points to the directory containing the migration SQL files.
var migrationsPath = "../migrations"

// setupTestPostgres spins up a temporary Postgres container, runs migrations, and returns a connected *sqlx.DB and a teardown function.
func setupTestPostgres(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Start a new Postgres container for testing
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Get the host and port to connect to the container
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Connect to the Postgres database using sqlx
	db, err := sqlx.Open("postgres", dsn)
	require.NoError(t, err)

	// Wait for the database to be ready to accept connections
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Run migrations using golang-migrate to set up the schema
	driver, err := migratepg.WithInstance(db.DB, &migratepg.Config{})
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	// Return the db connection and a teardown function to clean up resources
	teardown := func() {
		_ = db.Close()
		_ = container.Terminate(ctx)
	}
	return db, teardown
}

// TestPostgresAlbumRepository_Create tests only the Create method.
func TestPostgresAlbumRepository_Create(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

	album := Album{Title: "Where Did Our Love Go", Artist: "The Supremes", Price: 9.99, Year: 1964}
	err := repo.Create(album)
	require.NoError(t, err)

	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, albums, 1)
	require.Equal(t, "Where Did Our Love Go", albums[0].Title)
	require.Equal(t, "The Supremes", albums[0].Artist)
	require.Equal(t, 9.99, albums[0].Price)
	require.Equal(t, 1964, albums[0].Year)
	require.NotEmpty(t, albums[0].ID)
}

// TestPostgresAlbumRepository_GetAll tests only the GetAll method.
func TestPostgresAlbumRepository_GetAll(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970})
	_ = repo.Create(Album{Title: "Diana", Artist: "Diana Ross", Price: 2.0, Year: 1980})
	_ = repo.Create(Album{Title: "Sex Machine", Artist: "James Brown", Price: 3.0, Year: 1970})

	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, albums)
	require.Equal(t, "ABC", albums[0].Title)
	require.Equal(t, "Jackson 5", albums[0].Artist)
	require.Equal(t, 1.0, albums[0].Price)
	require.Equal(t, 1970, albums[0].Year)

	require.Equal(t, "Diana", albums[1].Title)
	require.Equal(t, "Diana Ross", albums[1].Artist)
	require.Equal(t, 2.0, albums[1].Price)
	require.Equal(t, 1980, albums[1].Year)

	require.Equal(t, "Sex Machine", albums[2].Title)
	require.Equal(t, "James Brown", albums[2].Artist)
	require.Equal(t, 3.0, albums[2].Price)
	require.Equal(t, 1970, albums[2].Year)
}

// TestPostgresAlbumRepository_GetByID tests only the GetByID method.
func TestPostgresAlbumRepository_GetByID(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970})
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
	require.Equal(t, id, got.ID)
}

// TestPostgresAlbumRepository_Delete tests only the Delete method.
func TestPostgresAlbumRepository_Delete(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970})
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
