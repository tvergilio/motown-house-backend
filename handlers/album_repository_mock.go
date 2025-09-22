package handlers

import (
	"os"
	"strconv"

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
		if strconv.Itoa(a.ID) == id {
			return a, nil
		}
	}
	return repository.Album{}, os.ErrNotExist
}

func (m *mockAlbumRepo) Create(album repository.Album) error {
	m.albums = append(m.albums, album)
	return nil
}

func (m *mockAlbumRepo) Delete(id string) error {
	for i, a := range m.albums {
		if strconv.Itoa(a.ID) == id {
			m.albums = append(m.albums[:i], m.albums[i+1:]...)
			return nil
		}
	}
	return os.ErrNotExist
}

func newTestHandler() *AlbumHandler {
	initial := []repository.Album{
		{ID: 101, Title: "Thriller", Artist: "Michael Jackson", Price: 42.99},
		{ID: 102, Title: "Lady Soul", Artist: "Aretha Franklin", Price: 35.50},
		{ID: 103, Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.00},
	}
	repo := &mockAlbumRepo{albums: initial}
	return NewAlbumHandler(repo)
}
