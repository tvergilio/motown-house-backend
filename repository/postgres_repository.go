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
	err := r.db.Select(&albums, "SELECT id, title, artist, price, year FROM albums")
	return albums, err
}

func (r *PostgresAlbumRepository) GetByID(id int) (Album, error) {
	var album Album
	err := r.db.Get(&album, "SELECT id, title, artist, price, year FROM albums WHERE id = $1", id)
	return album, err
}

func (r *PostgresAlbumRepository) Create(album Album) error {
	_, err := r.db.Exec(
		"INSERT INTO albums (title, artist, price, year) VALUES ($1, $2, $3, $4)",
		album.Title, album.Artist, album.Price, album.Year,
	)
	return err
}

func (r *PostgresAlbumRepository) Update(album Album) error {
	_, err := r.db.Exec(
		"UPDATE albums SET title = $1, artist = $2, price = $3, year = $4 WHERE id = $5",
		album.Title, album.Artist, album.Price, album.Year, album.ID,
	)
	return err
}

func (r *PostgresAlbumRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM albums WHERE id = $1", id)
	return err
}
