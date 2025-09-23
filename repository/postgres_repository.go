package repository

import (
	"github.com/jmoiron/sqlx"
)

type PostgresAlbumRepository struct {
	db *sqlx.DB
}

func NewPostgresAlbumRepository(db *sqlx.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

func (r *PostgresAlbumRepository) GetAll() ([]Album, error) {
	var albums []Album
	err := r.db.Select(&albums, "SELECT id, title, artist, price FROM albums")
	return albums, err
}

func (r *PostgresAlbumRepository) GetByID(id int) (Album, error) {
	var album Album
	err := r.db.Get(&album, "SELECT id, title, artist, price FROM albums WHERE id = $1", id)
	return album, err
}

func (r *PostgresAlbumRepository) Create(album Album) error {
	_, err := r.db.Exec(
		"INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)",
		album.Title, album.Artist, album.Price,
	)
	return err
}

func (r *PostgresAlbumRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM albums WHERE id = $1", id)
	return err
}
