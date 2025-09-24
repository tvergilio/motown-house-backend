package repository

type Album struct {
	ID     int     `db:"id" json:"id"`
	Title  string  `db:"title" json:"title"`
	Artist string  `db:"artist" json:"artist"`
	Price  float64 `db:"price" json:"price"`
	Year   int     `db:"year" json:"year"`
}

type AlbumRepository interface {
	GetAll() ([]Album, error)
	GetByID(id int) (Album, error)
	Create(album Album) error
	Delete(id int) error
}
