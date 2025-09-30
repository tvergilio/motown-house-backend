package repository

type Album struct {
	ID       string  `db:"id" json:"id"`
	Title    string  `db:"title" json:"title"`
	Artist   string  `db:"artist" json:"artist"`
	Price    float64 `db:"price" json:"price"`
	Year     int     `db:"year" json:"year"`
	ImageUrl string  `db:"image_url" json:"imageUrl"`
	Genre    string  `db:"genre" json:"genre"`
}

type AlbumRepository interface {
	GetAll() ([]Album, error)
	GetByID(id string) (Album, error)
	Create(album Album) error
	Delete(id string) error
	Update(album Album) error
}
