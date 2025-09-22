package main

import (
	"log"

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
		{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Title: "Giant Steps", Artist: "John Coltrane", Price: 63.99},
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
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

	router := gin.Default()
	router.GET("/albums", handler.GetAlbums)
	router.GET("/albums/:id", handler.GetAlbumByID)
	router.POST("/albums", handler.PostAlbums)
	router.DELETE("/albums/:id", handler.DeleteAlbum)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
