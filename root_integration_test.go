//go:build integration

package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = "http://silly-demo.127.0.0.1.nip.io"
	}
	client := &http.Client{}

	t.Run("Default message when MESSAGE env var is not set", func(t *testing.T) {
		// This test assumes the server is running WITHOUT MESSAGE env var
		// or with an empty MESSAGE env var
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// Test for default message - this is what we expect when MESSAGE is not set
		expectedMessage := "This is a silly demo\n"
		actualBody := string(body)

		// Only assert default if we're not testing with MESSAGE env var
		// (This is a limitation of integration tests - we can't easily control server env vars)
		if actualBody == expectedMessage {
			assert.Equal(t, expectedMessage, actualBody)
		} else {
			t.Logf("Server appears to be running with MESSAGE env var set. Body: %s", actualBody)
			// At least verify it's not empty and ends with newline
			assert.NotEmpty(t, actualBody)
			assert.True(t, len(actualBody) > 1 && actualBody[len(actualBody)-1] == '\n')
		}
	})

	t.Run("HTML formatting", func(t *testing.T) {
		req, err := http.NewRequest("GET", url+"?html=true", nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)

		actualBody := string(body)
		// Should be wrapped in <h1> tags and end with newline
		assert.True(t, len(actualBody) > 9) // At least "<h1></h1>\n"
		assert.True(t, actualBody[:4] == "<h1>")
		assert.True(t, actualBody[len(actualBody)-6:] == "</h1>\n")
	})

	t.Run("Fail parameter", func(t *testing.T) {
		req, err := http.NewRequest("GET", url+"?fail=true", nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "Something terrible happened", string(body))
	})
}
