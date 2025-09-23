package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/web-service-gin/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func Test_GetAlbums_StatusAndContent(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAlbums(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Thriller")
	assert.Contains(t, w.Body.String(), "Lady Soul")
	assert.Contains(t, w.Body.String(), "What's Going On")
}

func Test_GetAlbums_ValidJSON(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAlbums(c)

	var resp []repository.Album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "response should be valid JSON")
}

func Test_GetAlbums_ResponseIsArray(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAlbums(c)

	var resp []repository.Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.IsType(t, []repository.Album{}, resp, "response should be a slice of Album")
}

func Test_GetAlbums_CorrectAlbumCount(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAlbums(c)

	var resp []repository.Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 3, len(resp), "should return 3 albums")
}

func setupRouter(handler *AlbumHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/albums", handler.GetAlbums)
	r.GET("/albums/:id", handler.GetAlbumByID)
	r.POST("/albums", handler.PostAlbums)
	r.DELETE("/albums/:id", handler.DeleteAlbum)
	return r
}

func Test_GetAlbumByID_Found(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/albums/101", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var album repository.Album
	err := json.Unmarshal(w.Body.Bytes(), &album)
	assert.NoError(t, err)
	assert.Equal(t, "Thriller", album.Title)
	assert.Equal(t, "Michael Jackson", album.Artist)
	assert.Equal(t, 42.99, album.Price)
}

func Test_GetAlbumByID_NotFound(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/albums/999", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "album not found")
}

func Test_PostAlbums_Success(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	album := repository.Album{ID: 104, Title: "Bad", Artist: "Michael Jackson", Price: 29.99}
	jsonBytes, _ := json.Marshal(album)
	c.Request = httptest.NewRequest("POST", "/albums", io.NopCloser(bytes.NewReader(jsonBytes)))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostAlbums(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp repository.Album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, album.Title, resp.Title)
	assert.Equal(t, album.Artist, resp.Artist)
}

func Test_PostAlbums_InvalidJSON(t *testing.T) {
	handler := newTestHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/albums", io.NopCloser(bytes.NewReader([]byte("invalid json"))))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostAlbums(c)

	// Handler returns 400 on invalid JSON
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func Test_DeleteAlbum_Success(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/albums/101", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_DeleteAlbum_NotFound(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/albums/999", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "album not found")
}
