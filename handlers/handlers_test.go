package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func resetAlbums() {
	Albums = []Album{
		{ID: "101", Title: "Thriller", Artist: "Michael Jackson", Price: 42.99},
		{ID: "102", Title: "Lady Soul", Artist: "Aretha Franklin", Price: 35.50},
		{ID: "103", Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.00},
	}
}

func Test_GetAlbums_StatusAndContent(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetAlbums(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Thriller")
	assert.Contains(t, w.Body.String(), "Lady Soul")
	assert.Contains(t, w.Body.String(), "What's Going On")
}

func Test_GetAlbums_ValidJSON(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetAlbums(c)

	var resp []Album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "response should be valid JSON")
}

func Test_GetAlbums_ResponseIsArray(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetAlbums(c)

	var resp []Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.IsType(t, []Album{}, resp, "response should be a slice of Album")
}

func Test_GetAlbums_CorrectAlbumCount(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetAlbums(c)

	var resp []Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 3, len(resp), "should return 3 albums")
}

func Test_GetAlbumByID_Found(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "101"}}

	GetAlbumByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var album Album
	err := json.Unmarshal(w.Body.Bytes(), &album)
	assert.NoError(t, err)
	assert.Equal(t, "Thriller", album.Title)
	assert.Equal(t, "Michael Jackson", album.Artist)
	assert.Equal(t, 42.99, album.Price)
}

func Test_GetAlbumByID_NotFound(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	GetAlbumByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "album not found")
}

func Test_PostAlbums_Success(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	album := Album{ID: "104", Title: "Bad", Artist: "Michael Jackson", Price: 29.99}
	jsonBytes, _ := json.Marshal(album)
	c.Request = httptest.NewRequest("POST", "/albums", io.NopCloser(bytes.NewReader(jsonBytes)))
	c.Request.Header.Set("Content-Type", "application/json")

	PostAlbums(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp Album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, album.Title, resp.Title)
	assert.Equal(t, album.Artist, resp.Artist)
}

func Test_PostAlbums_InvalidJSON(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/albums", io.NopCloser(bytes.NewReader([]byte("invalid json"))))
	c.Request.Header.Set("Content-Type", "application/json")

	PostAlbums(c)

	// Handler returns 400 on invalid JSON
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Empty(t, w.Body.String())
}
