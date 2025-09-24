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
		{ID: 101, Title: "Thriller", Artist: "Michael Jackson", Price: 42.99, Year: 1982},
		{ID: 102, Title: "Lady Soul", Artist: "Aretha Franklin", Price: 35.50, Year: 1968},
		{ID: 103, Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.00, Year: 1971},
	}
	repo := &mockAlbumRepo{albums: initial}
	return NewAlbumHandler(repo)
}
