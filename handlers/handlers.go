package handlers

import (
	"net/http"
	"strconv"

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

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	idStr := c.Param("id")
	if _, err := strconv.Atoi(idStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	album, err := h.Repo.GetByID(idStr)
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
