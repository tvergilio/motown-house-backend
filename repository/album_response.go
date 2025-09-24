package repository

// AlbumResponse is what the /api/search endpoint will return to the frontend.
type AlbumResponse struct {
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	Price    float64 `json:"price"`
	Year     int     `json:"year"`
	Genre    string  `json:"genre"`
	ImageURL string  `json:"image_url"`
}
