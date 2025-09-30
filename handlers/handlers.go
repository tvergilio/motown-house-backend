package handlers

import (
	"fmt"
	"net/http"

	"example.com/web-service-gin/repository"
	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	Repo       repository.AlbumRepository
	ITunesRepo repository.ITunesRepository
}

func NewAlbumHandler(repo repository.AlbumRepository, itunesRepo repository.ITunesRepository) *AlbumHandler {
	return &AlbumHandler{
		Repo:       repo,
		ITunesRepo: itunesRepo,
	}
}

func (h *AlbumHandler) GetAlbums(c *gin.Context) {
	albums, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// AlbumIDUri is used for binding and validating the `id` URI parameter in routes like /albums/:id.
// This struct is specific to HTTP request handling and should not be used in the domain or repository layers.
type AlbumIDUri struct {
	ID string `uri:"id" binding:"required"`
}

// getAlbumIDFromUri extracts and validates the `id` URI parameter from the context.
// Returns the id as string and true if successful, otherwise writes a standardised error response and returns false.
func getAlbumIDFromUri(c *gin.Context) (string, bool) {
	var uri AlbumIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid album ID"})
		return "", false
	}
	return uri.ID, true
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id, ok := getAlbumIDFromUri(c)
	if !ok {
		return
	}
	album, err := h.Repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func (h *AlbumHandler) PostAlbums(c *gin.Context) {
	var newAlbum repository.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate required fields
	if newAlbum.ImageUrl == "" || newAlbum.Genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageUrl and genre are required and cannot be empty"})
		return
	}
	if err := h.Repo.Create(newAlbum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func (h *AlbumHandler) PutAlbum(c *gin.Context) {
	id, ok := getAlbumIDFromUri(c)
	if !ok {
		return
	}
	var updatedAlbum repository.Album
	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate required fields
	if updatedAlbum.ImageUrl == "" || updatedAlbum.Genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageUrl and genre are required and cannot be empty"})
		return
	}
	updatedAlbum.ID = id
	err := h.Repo.Update(updatedAlbum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedAlbum)
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id, ok := getAlbumIDFromUri(c)
	if !ok {
		return
	}
	err := h.Repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// SearchAlbums handles GET /api/search endpoint to search iTunes for albums
func (h *AlbumHandler) SearchAlbums(c *gin.Context) {
	// Get the search term from the query parameter
	term := c.Query("term")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "term query parameter is required"})
		return
	}

	// Use iTunes repository to search for albums
	searchResults, err := h.ITunesRepo.Search(term)
	if err != nil {
		message := fmt.Sprintf("failed to fetch data from iTunes API: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, searchResults)
}
