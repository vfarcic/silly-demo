//go:build unit

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
	// Save original environment variables
	originalMessage := os.Getenv("MESSAGE")
	originalVersion := os.Getenv("VERSION")
	originalFail := os.Getenv("FAIL")

	// Clean up after tests
	defer func() {
		if originalMessage != "" {
			os.Setenv("MESSAGE", originalMessage)
		} else {
			os.Unsetenv("MESSAGE")
		}
		if originalVersion != "" {
			os.Setenv("VERSION", originalVersion)
		} else {
			os.Unsetenv("VERSION")
		}
		if originalFail != "" {
			os.Setenv("FAIL", originalFail)
		} else {
			os.Unsetenv("FAIL")
		}
	}()

	// Clear environment variables for clean testing
	os.Unsetenv("MESSAGE")
	os.Unsetenv("VERSION")
	os.Unsetenv("FAIL")

	t.Run("No query parameters and no MESSAGE env var", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "This is a silly demo\n", w.Body.String())
	})

	t.Run("Query parameter 'fail' is present", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?fail=true", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Something terrible happened", w.Body.String())
	})

	t.Run("Environment variable FAIL is set", func(t *testing.T) {
		os.Setenv("FAIL", "true")
		defer os.Unsetenv("FAIL")

		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Something terrible happened", w.Body.String())
	})

	t.Run("Query parameter 'html' is present", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?html=true", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "<h1>This is a silly demo</h1>\n", w.Body.String())
	})

	t.Run("Environment variable MESSAGE is set", func(t *testing.T) {
		os.Setenv("MESSAGE", "Custom message from environment")
		defer os.Unsetenv("MESSAGE")

		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Custom message from environment\n", w.Body.String())
	})

	t.Run("Environment variable MESSAGE is set with HTML formatting", func(t *testing.T) {
		os.Setenv("MESSAGE", "Custom HTML message")
		defer os.Unsetenv("MESSAGE")

		req, _ := http.NewRequest("GET", "/?html=true", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "<h1>Custom HTML message</h1>\n", w.Body.String())
	})

	t.Run("Environment variable MESSAGE is empty string", func(t *testing.T) {
		os.Setenv("MESSAGE", "")
		defer os.Unsetenv("MESSAGE")

		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/", rootHandler)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "This is a silly demo\n", w.Body.String())
	})

	// Note: VERSION environment variable tests are commented out because
	// includeVersion is hardcoded to false in the current implementation
	// If you want to enable version functionality, uncomment these tests:

	// t.Run("Environment variable VERSION is set but includeVersion is false", func(t *testing.T) {
	// 	os.Setenv("VERSION", "1.0.0")
	// 	defer os.Unsetenv("VERSION")
	//
	// 	req, _ := http.NewRequest("GET", "/", nil)
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.GET("/", rootHandler)
	// 	router.ServeHTTP(w, req)
	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Equal(t, "This is a silly demo\n", w.Body.String())
	// })

	// t.Run("MESSAGE and VERSION both set with includeVersion true", func(t *testing.T) {
	// 	os.Setenv("MESSAGE", "Test app")
	// 	os.Setenv("VERSION", "2.0.0")
	// 	defer func() {
	// 		os.Unsetenv("MESSAGE")
	// 		os.Unsetenv("VERSION")
	// 	}()
	//
	// 	// This test would require modifying rootHandler to make includeVersion configurable
	// 	req, _ := http.NewRequest("GET", "/", nil)
	// 	w := httptest.NewRecorder()
	// 	router := gin.Default()
	// 	router.GET("/", rootHandler)
	// 	router.ServeHTTP(w, req)
	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Equal(t, "Test app version 2.0.0\n", w.Body.String())
	// })
}
