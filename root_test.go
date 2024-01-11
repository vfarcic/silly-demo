package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	// Test case 1: No query parameters
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/", rootHandler)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "This is a silly demo\n", w.Body.String())

	// Test case 2: Query parameter "fail" is present
	req, _ = http.NewRequest("GET", "/?fail=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Something terrible happened", w.Body.String())

	// Test case 3: Query parameter "html" is present
	req, _ = http.NewRequest("GET", "/?html=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<h1>This is a silly demo</h1>\n", w.Body.String())

	// Test case 4: Query parameter "html" and "fail" are present
	req, _ = http.NewRequest("GET", "/?html=true&fail=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Something terrible happened", w.Body.String())

	// Test case 5: Environment variable "VERSION" is set
	os.Setenv("VERSION", "1.0")
	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "This is a silly demo version 1.0\n", w.Body.String())
	os.Unsetenv("VERSION")

	// Test case 6: Environment variable "MESSAGE" is set
	os.Setenv("MESSAGE", "Custom message")
	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Custom message\n", w.Body.String())
	os.Unsetenv("MESSAGE")
}
