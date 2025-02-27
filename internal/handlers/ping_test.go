//go:build unit

package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test case 1: Valid URL
	t.Run("Valid URL", func(t *testing.T) {
		router := gin.Default()
		router.GET("/ping", PingHandler)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", "/ping?url="+server.URL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})

	// Test case 2: Empty URL with PING_URL environment variable
	t.Run("Empty URL with PING_URL environment variable", func(t *testing.T) {
		router := gin.Default()
		router.GET("/ping", PingHandler)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}))
		defer server.Close()

		os.Setenv("PING_URL", server.URL)
		defer os.Unsetenv("PING_URL")

		req, _ := http.NewRequest("GET", "/ping", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})

	// Test case 3: Empty URL without PING_URL environment variable
	t.Run("Empty URL without PING_URL environment variable", func(t *testing.T) {
		router := gin.Default()
		router.GET("/ping", PingHandler)

		req, _ := http.NewRequest("GET", "/ping", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "url is empty")
	})

	// Test case 4: Invalid URL
	t.Run("Invalid URL", func(t *testing.T) {
		router := gin.Default()
		router.GET("/ping", PingHandler)

		req, _ := http.NewRequest("GET", "/ping?url=invalid-url", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "unsupported protocol scheme")
	})

}
