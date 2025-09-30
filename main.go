package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tvergilio/motown-house-backend/db"
	"github.com/tvergilio/motown-house-backend/handlers"
	"github.com/tvergilio/motown-house-backend/repository"
)

func seedAlbums(repo repository.AlbumRepository) {
	albums, err := repo.GetAll()
	if err != nil {
		log.Printf("seedAlbums: failed to get albums: %v", err)
		return
	}
	if len(albums) > 0 {
		log.Printf("seedAlbums: albums table already seeded (%d records)", len(albums))
		return
	}
	initialAlbums := []repository.Album{
		{Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.99, Year: 1971, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music112/v4/76/36/2d/76362d74-cb7a-8ef9-104e-cde1d858e9a9/20UMGIM95279.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"},
		{Title: "Songs in the Key of Life", Artist: "Stevie Wonder", Price: 42.50, Year: 1976, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music118/v4/eb/1f/12/eb1f12ec-474c-63aa-43af-09282f423b9d/00602537004737.rgb.jpg/100x100bb.jpg", Genre: "Motown"},
		{Title: "Diana", Artist: "Diana Ross", Price: 28.75, Year: 1980, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music115/v4/aa/87/1c/aa871c20-95be-38bd-97e3-ecfeb8ec404b/15UMGIM06551.rgb.jpg/100x100bb.jpg", Genre: "R&B/Soul"},
		{Title: "Sex Machine", Artist: "James Brown", Price: 3.0, Year: 1970, ImageUrl: "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/17/8b/05/178b05de-5855-0136-9827-a0e8a6ccf3db/00602547021656.rgb.jpg/100x100bb.jpg", Genre: "Soul"},
	}
	for _, album := range initialAlbums {
		if err := repo.Create(album); err != nil {
			log.Printf("seedAlbums: failed to create album %+v: %v", album, err)
		}
	}
	log.Printf("seedAlbums: seeded %d albums", len(initialAlbums))
}

func main() {
	_ = godotenv.Load()

	dbConn, err := db.ConnectFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	// Handle error from Close
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Select repository implementation based on database backend
	var repo repository.AlbumRepository
	switch dbConn.Backend {
	case "postgres":
		log.Printf("Using Postgres backend")
		repo = repository.NewPostgresAlbumRepository(dbConn.PostgresDB)
	case "cassandra":
		log.Printf("Using Cassandra backend")
		repo = repository.NewCassandraAlbumRepository(dbConn.CassandraDB)
	default:
		log.Fatalf("Unsupported database backend: %s", dbConn.Backend)
	}

	itunesRepo := repository.NewITunesRepository()
	seedAlbums(repo)
	handler := handlers.NewAlbumHandler(repo, itunesRepo)

	r := gin.Default()

	// Enable CORS for frontend (supports both localhost and Docker network)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // Local development
			"http://frontend:3000",  // Docker network
			"http://127.0.0.1:3000", // Alternative localhost
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/albums", handler.GetAlbums)
	r.GET("/albums/:id", handler.GetAlbumByID)
	r.POST("/albums", handler.PostAlbums)
	r.DELETE("/albums/:id", handler.DeleteAlbum)
	r.PUT("/albums/:id", handler.PutAlbum)
	r.GET("/api/search", handler.SearchAlbums)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
