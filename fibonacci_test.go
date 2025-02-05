//go:build unit

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFibonacciHandler(t *testing.T) {
	// Test case 1: Valid number
	req, _ := http.NewRequest("GET", "/fibonacci?number=5", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/fibonacci", fibonacciHandler)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(5), response["number"])
	assert.Equal(t, float64(5), response["fibonacci"])

	// Test case 2: Invalid input
	req, _ = http.NewRequest("GET", "/fibonacci?number=abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"abc\": invalid syntax", response["error"])

	// Test case 3: Missing number parameter
	req, _ = http.NewRequest("GET", "/fibonacci", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"\": invalid syntax", response["error"])
}
