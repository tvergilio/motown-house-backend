package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ITunesSearchResponse maps the top-level response from iTunes.
type ITunesSearchResponse struct {
	ResultCount int           `json:"resultCount"`
	Results     []ITunesAlbum `json:"results"`
}

// ITunesAlbum maps one album entry from iTunes.
type ITunesAlbum struct {
	ArtistName       string  `json:"artistName"`
	CollectionName   string  `json:"collectionName"`
	CollectionPrice  float64 `json:"collectionPrice"`
	ReleaseDate      string  `json:"releaseDate"`
	PrimaryGenreName string  `json:"primaryGenreName"`
	ArtworkUrl100    string  `json:"artworkUrl100"`
}

// ITunesRepository interface for searching iTunes API
type ITunesRepository interface {
	Search(term string) ([]AlbumResponse, error)
}

// ITunesRepositoryImpl implements ITunesRepository
type ITunesRepositoryImpl struct {
	baseURL string
}

// NewITunesRepository creates a new iTunes repository
func NewITunesRepository() ITunesRepository {
	return &ITunesRepositoryImpl{
		baseURL: "https://itunes.apple.com/search",
	}
}

// Search searches iTunes API for albums matching the given term
func (r *ITunesRepositoryImpl) Search(term string) ([]AlbumResponse, error) {
	if term == "" {
		return nil, fmt.Errorf("search term cannot be empty")
	}

	// URL encode the search term to handle spaces and special characters
	encodedTerm := url.QueryEscape(term)

	// Make a request to iTunes Search API
	itunesURL := fmt.Sprintf("%s?term=%s&entity=album", r.baseURL, encodedTerm)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(itunesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from iTunes API: %w", err)
	}
	defer resp.Body.Close()

	// Check if the iTunes API request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes API returned status code: %d", resp.StatusCode)
	}

	// Parse JSON response
	var itunesResponse ITunesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&itunesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse iTunes API response: %w", err)
	}

	// Convert iTunes albums to our AlbumResponse format
	searchResults := make([]AlbumResponse, 0, len(itunesResponse.Results))
	for _, itunesAlbum := range itunesResponse.Results {
		// Extract year from release date (format: "1970-09-01T07:00:00Z")
		year := 0
		if itunesAlbum.ReleaseDate != "" {
			if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", itunesAlbum.ReleaseDate); err == nil {
				year = parsedTime.Year()
			}
		}

		albumResponse := AlbumResponse{
			Title:    itunesAlbum.CollectionName,
			Artist:   itunesAlbum.ArtistName,
			Price:    itunesAlbum.CollectionPrice,
			Year:     year,
			Genre:    itunesAlbum.PrimaryGenreName,
			ImageURL: itunesAlbum.ArtworkUrl100,
		}
		searchResults = append(searchResults, albumResponse)
	}

	return searchResults, nil
}
