package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestITunesRepository_Search_Integration tests the actual iTunes API integration
func TestITunesRepository_Search_Integration(t *testing.T) {
	// Skip this test in CI/CD or when running without internet
	if testing.Short() {
		t.Skip("Skipping iTunes API integration test in short mode")
	}

	repo := NewITunesRepository()

	// Test with a well-known album that should exist in iTunes
	results, err := repo.Search("Michael Jackson Thriller")
	require.NoError(t, err)
	require.NotEmpty(t, results, "Should find at least one result for Michael Jackson Thriller")

	// Verify the structure of the first result
	firstResult := results[0]
	require.NotEmpty(t, firstResult.Title, "Title should not be empty")
	require.NotEmpty(t, firstResult.Artist, "Artist should not be empty")
	require.NotEmpty(t, firstResult.Genre, "Genre should not be empty")
	require.NotEmpty(t, firstResult.ImageURL, "ImageURL should not be empty")
	require.Greater(t, firstResult.Year, 0, "Year should be greater than 0")
	require.Greater(t, firstResult.Price, 0.0, "Price should be greater than 0")

	// Test that we can find Michael Jackson in the results
	foundMichaelJackson := false
	for _, result := range results {
		if result.Artist == "Michael Jackson" {
			foundMichaelJackson = true
			break
		}
	}
	require.True(t, foundMichaelJackson, "Should find at least one Michael Jackson album")
}

// TestITunesRepository_Search_EmptyTerm tests error handling for an empty search term
func TestITunesRepository_Search_EmptyTerm(t *testing.T) {
	repo := NewITunesRepository()

	results, err := repo.Search("")
	require.Error(t, err)
	require.Nil(t, results)
	require.Contains(t, err.Error(), "search term cannot be empty")
}

// TestITunesRepository_Search_UncommonTerm tests search with uncommon term
func TestITunesRepository_Search_UncommonTerm(t *testing.T) {
	// Skip this test in CI/CD or when running without an internet connection
	if testing.Short() {
		t.Skip("Skipping iTunes API integration test in short mode")
	}

	repo := NewITunesRepository()

	// Test with a very uncommon search term that likely won't return results
	_, err := repo.Search("xyzabc123nonexistentalbum")
	require.NoError(t, err, "Should not error even with no results")
	// Results could be empty or contain unexpected matches, both are valid
}

// TestITunesRepository_Search_SpecialCharacters tests search with special characters
func TestITunesRepository_Search_SpecialCharacters(t *testing.T) {
	// Skip this test in CI/CD or when running without an internet connection
	if testing.Short() {
		t.Skip("Skipping iTunes API integration test in short mode")
	}

	repo := NewITunesRepository()

	// Test with special characters that need URL encoding
	_, err := repo.Search("Diana Ross & The Supremes")
	require.NoError(t, err)
	// Should handle URL encoding properly without errors
}
