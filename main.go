package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"example.com/web-service-gin/db"
	"example.com/web-service-gin/handlers"
	"example.com/web-service-gin/repository"
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
		{Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.99, Year: 1971},
		{Title: "Songs in the Key of Life", Artist: "Stevie Wonder", Price: 42.50, Year: 1976},
		{Title: "Diana", Artist: "Diana Ross", Price: 28.75, Year: 1980},
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
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	// Handle error from Close
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	repo := repository.NewPostgresAlbumRepository(database)
	seedAlbums(repo)
	handler := handlers.NewAlbumHandler(repo)

	r := gin.Default()

	// Enable CORS for frontend (e.g., http://localhost:9002)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/albums", handler.GetAlbums)
	r.GET("/albums/:id", handler.GetAlbumByID)
	r.POST("/albums", handler.PostAlbums)
	r.DELETE("/albums/:id", handler.DeleteAlbum)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
