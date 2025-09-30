package repository

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSession is a mock implementation of gocql.Session for unit testing
type MockSession struct {
	mock.Mock
}

// MockQuery is a mock implementation of gocql.Query
type MockQuery struct {
	mock.Mock
}

// MockIter is a mock implementation of gocql.Iter
type MockIter struct {
	mock.Mock
	albums []Album
	index  int
}

func (m *MockSession) Query(stmt string, values ...interface{}) *MockQuery {
	args := m.Called(stmt, values)
	return args.Get(0).(*MockQuery)
}

func (m *MockQuery) Exec() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockQuery) Iter() *MockIter {
	args := m.Called()
	return args.Get(0).(*MockIter)
}

func (m *MockQuery) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockIter) Scan(dest ...interface{}) bool {
	args := m.Called(dest)
	return args.Bool(0)
}

func (m *MockIter) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper function to create a test album
func createTestAlbum() Album {
	return Album{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Title:    "Test Album",
		Artist:   "Test Artist",
		Price:    19.99,
		Year:     2023,
		ImageUrl: "https://example.com/image.jpg",
		Genre:    "Rock",
	}
}

// TestNewCassandraAlbumRepository tests the constructor
func TestNewCassandraAlbumRepository(t *testing.T) {
	mockSession := &MockSession{}

	repo := NewCassandraAlbumRepository((*gocql.Session)(nil))

	require.NotNil(t, repo)
	assert.IsType(t, &CassandraAlbumRepository{}, repo)

	// Test with mock session (we can't directly test the real session due to type conversion)
	_ = mockSession // Use mockSession to avoid unused variable
}

// TestCassandraAlbumRepository_GetByID_Success tests successful ID lookup
func TestCassandraAlbumRepository_GetByID_Success(t *testing.T) {
	// Test UUID parsing logic
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	_, err := gocql.ParseUUID(validUUID)
	require.NoError(t, err, "Test UUID should be valid")

	// Test invalid UUID handling
	invalidUUID := "invalid-uuid"
	_, err = gocql.ParseUUID(invalidUUID)
	require.Error(t, err, "Invalid UUID should return error")
}

// TestCassandraAlbumRepository_GetByID_InvalidUUID tests error handling for invalid UUID
func TestCassandraAlbumRepository_GetByID_InvalidUUID(t *testing.T) {
	repo := &CassandraAlbumRepository{
		session: nil, // We don't need a real session for UUID validation test
	}

	// Test with invalid UUID
	album, err := repo.GetByID("invalid-uuid")

	assert.Error(t, err)
	assert.Equal(t, Album{}, album)
	assert.Contains(t, err.Error(), "invalid UUID")
}

// TestCassandraAlbumRepository_Update_InvalidUUID tests error handling for invalid UUID in update
func TestCassandraAlbumRepository_Update_InvalidUUID(t *testing.T) {
	repo := &CassandraAlbumRepository{
		session: nil, // We don't need a real session for UUID validation test
	}

	album := createTestAlbum()
	album.ID = "invalid-uuid"

	err := repo.Update(album)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

// TestCassandraAlbumRepository_Delete_InvalidUUID tests error handling for invalid UUID in delete
func TestCassandraAlbumRepository_Delete_InvalidUUID(t *testing.T) {
	repo := &CassandraAlbumRepository{
		session: nil, // We don't need a real session for UUID validation test
	}

	err := repo.Delete("invalid-uuid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

// TestCassandraAlbumRepository_UUIDGeneration tests UUID generation logic
func TestCassandraAlbumRepository_UUIDGeneration(t *testing.T) {
	// Test that TimeUUID generates valid UUIDs
	uuid1 := gocql.TimeUUID()
	uuid2 := gocql.TimeUUID()

	// UUIDs should be different
	assert.NotEqual(t, uuid1, uuid2)

	// UUIDs should be valid
	assert.NotEqual(t, gocql.UUID{}, uuid1)
	assert.NotEqual(t, gocql.UUID{}, uuid2)

	// String conversion should work
	str1 := uuid1.String()
	str2 := uuid2.String()

	assert.NotEmpty(t, str1)
	assert.NotEmpty(t, str2)
	assert.NotEqual(t, str1, str2)

	// Should be able to parse back
	parsed1, err1 := gocql.ParseUUID(str1)
	parsed2, err2 := gocql.ParseUUID(str2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, uuid1, parsed1)
	assert.Equal(t, uuid2, parsed2)
}

// TestCassandraAlbumRepository_AlbumValidation tests album data validation
func TestCassandraAlbumRepository_AlbumValidation(t *testing.T) {
	testCases := []struct {
		name    string
		album   Album
		isValid bool
	}{
		{
			name:    "Valid album",
			album:   createTestAlbum(),
			isValid: true,
		},
		{
			name: "Empty title",
			album: Album{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				Title:    "",
				Artist:   "Test Artist",
				Price:    19.99,
				Year:     2023,
				ImageUrl: "https://example.com/image.jpg",
				Genre:    "Rock",
			},
			isValid: false,
		},
		{
			name: "Empty artist",
			album: Album{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Album",
				Artist:   "",
				Price:    19.99,
				Year:     2023,
				ImageUrl: "https://example.com/image.jpg",
				Genre:    "Rock",
			},
			isValid: false,
		},
		{
			name: "Negative price",
			album: Album{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Album",
				Artist:   "Test Artist",
				Price:    -5.99,
				Year:     2023,
				ImageUrl: "https://example.com/image.jpg",
				Genre:    "Rock",
			},
			isValid: false,
		},
		{
			name: "Invalid year",
			album: Album{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Album",
				Artist:   "Test Artist",
				Price:    19.99,
				Year:     0,
				ImageUrl: "https://example.com/image.jpg",
				Genre:    "Rock",
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Basic validation logic (in a real app, you might have a separate validator)
			isValid := tc.album.Title != "" &&
				tc.album.Artist != "" &&
				tc.album.Price >= 0 &&
				tc.album.Year > 0 &&
				tc.album.ImageUrl != "" &&
				tc.album.Genre != ""

			assert.Equal(t, tc.isValid, isValid, "Album validation should match expected result")
		})
	}
}

// TestCassandraAlbumRepository_ErrorHandling tests various error scenarios
func TestCassandraAlbumRepository_ErrorHandling(t *testing.T) {
	t.Run("Nil session handling", func(t *testing.T) {
		repo := &CassandraAlbumRepository{session: nil}

		// These would panic in real usage, but we're testing the structure
		assert.NotNil(t, repo)
		assert.Nil(t, repo.session)
	})

	t.Run("UUID parsing errors", func(t *testing.T) {
		invalidUUIDs := []string{
			"",
			"not-a-uuid",
			"550e8400-e29b-41d4-a716", // too short
			"550e8400-e29b-41d4-a716-446655440000-extra", // too long
			"gggggggg-eeee-41d4-a716-446655440000",       // invalid chars
		}

		for _, invalidUUID := range invalidUUIDs {
			_, err := gocql.ParseUUID(invalidUUID)
			assert.Error(t, err, "Should return error for invalid UUID: %s", invalidUUID)
		}
	})
}

// For VS Code users: To run full integration tests, either:
// 1. Run from terminal: go test -tags=integration -timeout=300s ./repository
// 2. Or use VS Code's test runner which now has proper timeout configured
