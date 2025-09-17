package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_getAlbums_StatusAndContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Blue Train")
	assert.Contains(t, w.Body.String(), "Jeru")
	assert.Contains(t, w.Body.String(), "Sarah Vaughan and Clifford Brown")
}

func Test_getAlbums_ValidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "response should be valid JSON")
}

func Test_getAlbums_ResponseIsArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.IsType(t, []album{}, resp, "response should be a slice of album")
}

func Test_getAlbums_CorrectAlbumCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAlbums(c)

	var resp []album
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 3, len(resp), "should return 3 albums")
}
