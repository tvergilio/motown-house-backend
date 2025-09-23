package handlers

import (
	"net/http"

	"example.com/web-service-gin/repository"
	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	Repo repository.AlbumRepository
}

func NewAlbumHandler(repo repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{Repo: repo}
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
	ID int `uri:"id" binding:"required"`
}

// getAlbumIDFromUri extracts and validates the `id` URI parameter from the context.
// Returns the id as int and true if successful, otherwise writes an error response and returns false.
func getAlbumIDFromUri(c *gin.Context) (int, bool) {
	var uri AlbumIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
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
	if err := h.Repo.Create(newAlbum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
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
