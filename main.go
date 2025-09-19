package main

import (
	"example.com/web-service-gin/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", handlers.GetAlbums)
	router.GET("/albums/:id", handlers.GetAlbumByID)
	router.POST("/albums", handlers.PostAlbums)

	router.Run("localhost:8080")
}
