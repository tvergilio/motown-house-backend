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
	r.PUT("/albums/:id", handler.PutAlbum)
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
	assert.Equal(t, 1982, album.Year)
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

	album := repository.Album{ID: 104, Title: "Bad", Artist: "Michael Jackson", Price: 29.99, Year: 1987}
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

func Test_PutAlbums_Success(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	updatedAlbum := repository.Album{
		ID:     101,
		Title:  "Thriller 25",
		Artist: "Michael Jackson",
		Price:  45.99,
		Year:   1982,
	}
	jsonBytes, _ := json.Marshal(updatedAlbum)
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp repository.Album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, updatedAlbum.Title, resp.Title)
	assert.Equal(t, updatedAlbum.Price, resp.Price)
}

func Test_PutAlbums_InvalidJSON(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_PutAlbums_NonExistentID(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	updatedAlbum := repository.Album{ID: 999, Title: "Ghost Album", Artist: "Nobody", Price: 10.0, Year: 2000}
	jsonBytes, _ := json.Marshal(updatedAlbum)
	req := httptest.NewRequest("PUT", "/albums/999", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_PutAlbums_IDMismatch(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	updatedAlbum := repository.Album{ID: 102, Title: "Mismatch", Artist: "Test", Price: 20.0, Year: 2020}
	jsonBytes, _ := json.Marshal(updatedAlbum)
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp repository.Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 101, resp.ID) // Should use URI ID
}

func Test_PutAlbums_MissingFields(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	partial := map[string]interface{}{"Title": "Partial"}
	jsonBytes, _ := json.Marshal(partial)
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp repository.Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Partial", resp.Title)
}

func Test_PutAlbums_InvalidFieldTypes(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	body := `{"title":"Bad Type","artist":"Test","price":"not a float","year":"not an int"}`
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_PutAlbums_NegativePriceOrYear(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	updatedAlbum := repository.Album{ID: 101, Title: "Negative", Artist: "Test", Price: -10.0, Year: -1980}
	jsonBytes, _ := json.Marshal(updatedAlbum)
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp repository.Album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Price < 0)
	assert.True(t, resp.Year < 0)
}

func Test_PutAlbums_EmptyBody(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_PutAlbums_ExtraFields(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	body := `{"title":"Extra","artist":"Test","price":10.0,"year":2020,"extra":"field"}`
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_PutAlbums_LargePayload(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	largeTitle := make([]byte, 10000)
	for i := range largeTitle {
		largeTitle[i] = 'A'
	}
	updatedAlbum := repository.Album{ID: 101, Title: string(largeTitle), Artist: "Test", Price: 10.0, Year: 2020}
	jsonBytes, _ := json.Marshal(updatedAlbum)
	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_PutAlbums_MalformedJSON(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/albums/101", bytes.NewReader([]byte("{title:BadJson")))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
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

func Test_GetAlbumByID_InvalidID(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/albums/abc", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid album ID")
}

func Test_DeleteAlbum_InvalidID(t *testing.T) {
	handler := newTestHandler()
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/albums/abc", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid album ID")
}
