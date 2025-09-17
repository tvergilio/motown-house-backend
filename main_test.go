package main

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
	// setup code here
	gin.SetMode(gin.TestMode)
	code := m.Run() // run tests

	// teardown code here

	os.Exit(code)
}

func resetAlbums() {
	albums = []album{
		{ID: "101", Title: "Thriller", Artist: "Michael Jackson", Price: 42.99},
		{ID: "102", Title: "Lady Soul", Artist: "Aretha Franklin", Price: 35.50},
		{ID: "103", Title: "What's Going On", Artist: "Marvin Gaye", Price: 39.00},
	}
}

func Test_getAlbums_StatusAndContent(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Thriller")
	assert.Contains(t, w.Body.String(), "Lady Soul")
	assert.Contains(t, w.Body.String(), "What's Going On")
}

func Test_getAlbums_ValidJSON(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "response should be valid JSON")
}

func Test_getAlbums_ResponseIsArray(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.IsType(t, []album{}, resp, "response should be a slice of album")
}

func Test_getAlbums_CorrectAlbumCount(t *testing.T) {
	resetAlbums()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 3, len(resp), "should return 3 albums")
}
