package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"example.com/web-service-gin/config"
	"example.com/web-service-gin/db"
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
var migrationsPath = "../migrations/postgres"

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
				"POSTGRES_HOST_AUTH_METHOD": "trust", // Allow connections without password
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
	dsn := fmt.Sprintf("postgres://postgres@%s:%s/postgres?sslmode=disable", host, port.Port())

	// Create a test config using the dynamic container details
	testConfig := &config.Config{
		DBBackend:   "postgres",
		PostgresURL: dsn,
	}

	// Connect to the Postgres database using the config system
	database, err := db.Connect(testConfig)
	require.NoError(t, err)

	// Wait for the database to be ready to accept connections
	for i := 0; i < 10; i++ {
		if err := database.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Run migrations using golang-migrate to set up the schema
	driver, err := migratepg.WithInstance(database.DB, &migratepg.Config{})
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
		_ = database.Close()
		_ = container.Terminate(ctx)
	}
	return database, teardown
}

// TestPostgresAlbumRepository_Create tests only the Create method.
func TestPostgresAlbumRepository_Create(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

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

// TestPostgresAlbumRepository_GetAll tests only the GetAll method.
func TestPostgresAlbumRepository_GetAll(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

	_ = repo.Create(Album{Title: "ABC", Artist: "Jackson 5", Price: 1.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	_ = repo.Create(Album{Title: "Diana", Artist: "Diana Ross", Price: 2.0, Year: 1980, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/aa/87/1c/aa871c20-95be-38bd-97e3-ecfeb8ec404b/15UMGIM06551.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"})
	_ = repo.Create(Album{Title: "Sex Machine", Artist: "James Brown", Price: 3.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/17/8b/05/178b05de-5855-0136-9827-a0e8a6ccf3db/00602547021656.rgb.jpg/100x100bb.jpg", Genre: "Soul"})

	albums, err := repo.GetAll()
	require.NoError(t, err)
	require.NotEmpty(t, albums)
	require.Equal(t, "ABC", albums[0].Title)
	require.Equal(t, "Jackson 5", albums[0].Artist)
	require.Equal(t, 1.0, albums[0].Price)
	require.Equal(t, 1970, albums[0].Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music211/v4/cb/38/70/cb3870c2-1a9b-9310-e218-9d0f5a5e98f5/06UMGIM05267.rgb.jpg/100x100bb.jpg", albums[0].ImageUrl)
	require.Equal(t, "R&B/Soul", albums[0].Genre)
	require.NotEmpty(t, albums[0].ID)

	require.Equal(t, "Diana", albums[1].Title)
	require.Equal(t, "Diana Ross", albums[1].Artist)
	require.Equal(t, 2.0, albums[1].Price)
	require.Equal(t, 1980, albums[1].Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/aa/87/1c/aa871c20-95be-38bd-97e3-ecfeb8ec404b/15UMGIM06551.rgb.jpg/100x100bb.jpg", albums[1].ImageUrl)
	require.Equal(t, "R&B/Soul", albums[1].Genre)
	require.NotEmpty(t, albums[1].ID)

	require.Equal(t, "Sex Machine", albums[2].Title)
	require.Equal(t, "James Brown", albums[2].Artist)
	require.Equal(t, 3.0, albums[2].Price)
	require.Equal(t, 1970, albums[2].Year)
	require.Equal(t, "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/17/8b/05/178b05de-5855-0136-9827-a0e8a6ccf3db/00602547021656.rgb.jpg/100x100bb.jpg", albums[2].ImageUrl)
	require.Equal(t, "Soul", albums[2].Genre)
	require.NotEmpty(t, albums[2].ID)
}

// TestPostgresAlbumRepository_GetByID tests only the GetByID method.
func TestPostgresAlbumRepository_GetByID(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

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

// TestPostgresAlbumRepository_Update tests only the Update method.
func TestPostgresAlbumRepository_Update(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

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

// TestPostgresAlbumRepository_Delete tests only the Delete method.
func TestPostgresAlbumRepository_Delete(t *testing.T) {
	db, teardown := setupTestPostgres(t)
	defer teardown()
	repo := NewPostgresAlbumRepository(db)

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
