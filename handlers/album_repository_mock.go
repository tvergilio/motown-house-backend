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

func (m *mockAlbumRepo) GetByID(id int) (repository.Album, error) {
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

func (m *mockAlbumRepo) Delete(id int) error {
	for i, a := range m.albums {
		if a.ID == id {
			m.albums = append(m.albums[:i], m.albums[i+1:]...)
			return nil
		}
	}
	return os.ErrNotExist
}

func newTestHandler() *AlbumHandler {
	initial := []repository.Album{
		{ID: 101, Title: "Thriller", Artist: "Michael Jackson", Price: 42.99, Year: 1982, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/32/4f/fd/324ffda2-9e51-8f6a-0c2d-c6fd2b41ac55/074643811224.jpg/100x100bb.jpg", Genre: "Pop"},
		{ID: 102, Title: "Lady Soul", Artist: "Aretha Franklin", Price: 35.50, Year: 1968, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/e3/a7/ac/e3a7aca0-48e1-0882-8e25-3d68f7ba3a72/603497896646.jpg/100x100bb.jpg", Genre: "R&B/Soul"},
		{ID: 103, Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.00, Year: 1971, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music112/v4/76/36/2d/76362d74-cb7a-8ef9-104e-cde1d858e9a9/20UMGIM95279.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"},
	}
	repo := &mockAlbumRepo{albums: initial}
	return NewAlbumHandler(repo)
}
