package handlers

import (
	"encoding/json"
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
