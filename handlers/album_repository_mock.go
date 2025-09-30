package handlers

import (
	"os"

	"example.com/web-service-gin/repository"
)

// In-memory mock implementation of AlbumRepository for testing

type mockAlbumRepo struct {
	albums []repository.Album
}

func (m *mockAlbumRepo) GetAll() ([]repository.Album, error) {
	return m.albums, nil
}

func (m *mockAlbumRepo) GetByID(id string) (repository.Album, error) {
	for _, a := range m.albums {
		if a.ID == id {
			return a, nil
		}
	}
	return repository.Album{}, os.ErrNotExist
}

func (m *mockAlbumRepo) Create(album repository.Album) error {
	m.albums = append(m.albums, album)
	return nil
}

func (m *mockAlbumRepo) Update(album repository.Album) error {
	for i, a := range m.albums {
		if a.ID == album.ID {
			m.albums[i] = album
			return nil // Success
		}
	}
	return os.ErrNotExist // Return error if album not found
}

func (m *mockAlbumRepo) Delete(id string) error {
	for i, a := range m.albums {
		if a.ID == id {
			m.albums = append(m.albums[:i], m.albums[i+1:]...)
			return nil
		}
	}
	return os.ErrNotExist
}

// Mock iTunes repository for testing
type mockITunesRepo struct{}

func (m *mockITunesRepo) Search(term string) ([]repository.AlbumResponse, error) {
	// Return mock search results for testing
	return []repository.AlbumResponse{
		{
			Title:    "Thriller",
			Artist:   "Michael Jackson",
			Price:    9.99,
			Year:     1982,
			Genre:    "Pop",
			ImageURL: "https://example.com/thriller.jpg",
		},
		{
			Title:    "Bad",
			Artist:   "Michael Jackson",
			Price:    8.99,
			Year:     1987,
			Genre:    "Pop",
			ImageURL: "https://example.com/bad.jpg",
		},
	}, nil
}

// newTestHandler creates a test handler with mock repositories and pre-seeded data
func newTestHandler() *AlbumHandler {
	mockRepo := &mockAlbumRepo{
		albums: []repository.Album{
			{ID: "1", Title: "Thriller", Artist: "Michael Jackson", Price: 25.99, Year: 1982, ImageUrl: "https://example.com/thriller.jpg", Genre: "Pop"},
			{ID: "2", Title: "Songs in the Key of Life", Artist: "Stevie Wonder", Price: 42.50, Year: 1976, ImageUrl: "https://example.com/songs.jpg", Genre: "Motown"},
			{ID: "101", Title: "Thriller", Artist: "Michael Jackson", Price: 42.99, Year: 1982, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/32/4f/fd/324ffda2-9e51-8f6a-0c2d-c6fd2b41ac55/074643811224.jpg/100x100bb.jpg", Genre: "Pop"},
		},
	}
	mockITunesRepo := &mockITunesRepo{}
	return NewAlbumHandler(mockRepo, mockITunesRepo)
}
