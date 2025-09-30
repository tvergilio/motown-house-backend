package repository

import (
	"github.com/gocql/gocql"
)

type CassandraAlbumRepository struct {
	session *gocql.Session
}

func NewCassandraAlbumRepository(session *gocql.Session) *CassandraAlbumRepository {
	return &CassandraAlbumRepository{session: session}
}

func (r *CassandraAlbumRepository) GetAll() ([]Album, error) {
	var albums []Album

	iter := r.session.Query("SELECT id, title, artist, price, year, image_url, genre FROM albums").Iter()
	defer iter.Close()

	var album Album
	var cassandraID gocql.UUID
	for iter.Scan(&cassandraID, &album.Title, &album.Artist, &album.Price, &album.Year, &album.ImageUrl, &album.Genre) {
		album.ID = cassandraID.String() // Convert UUID to string
		albums = append(albums, album)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (r *CassandraAlbumRepository) GetByID(id string) (Album, error) {
	var album Album
	var cassandraID gocql.UUID

	// Parse the string ID back to UUID
	parsedUUID, err := gocql.ParseUUID(id)
	if err != nil {
		return Album{}, err
	}

	err = r.session.Query(
		"SELECT id, title, artist, price, year, image_url, genre FROM albums WHERE id = ? LIMIT 1",
		parsedUUID,
	).Scan(&cassandraID, &album.Title, &album.Artist, &album.Price, &album.Year, &album.ImageUrl, &album.Genre)

	if err != nil {
		return Album{}, err
	}

	album.ID = cassandraID.String()
	return album, nil
}

func (r *CassandraAlbumRepository) Create(album Album) error {
	// Generate a new UUID for the album
	albumID := gocql.TimeUUID()

	err := r.session.Query(
		"INSERT INTO albums (id, title, artist, price, year, image_url, genre) VALUES (?, ?, ?, ?, ?, ?, ?)",
		albumID, album.Title, album.Artist, album.Price, album.Year, album.ImageUrl, album.Genre,
	).Exec()

	return err
}

func (r *CassandraAlbumRepository) Update(album Album) error {
	// Parse the string ID back to UUID
	parsedUUID, err := gocql.ParseUUID(album.ID)
	if err != nil {
		return err
	}

	err = r.session.Query(
		"UPDATE albums SET title = ?, artist = ?, price = ?, year = ?, image_url = ?, genre = ? WHERE id = ?",
		album.Title, album.Artist, album.Price, album.Year, album.ImageUrl, album.Genre, parsedUUID,
	).Exec()

	return err
}

func (r *CassandraAlbumRepository) Delete(id string) error {
	// Parse the string ID back to UUID
	parsedUUID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	err = r.session.Query(
		"DELETE FROM albums WHERE id = ?",
		parsedUUID,
	).Exec()

	return err
}
