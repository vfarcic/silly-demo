//go:build integration

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	// Test case 1: No query parameters
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = "http://silly-demo.127.0.0.1.nip.io"
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer res.Body.Close()
	w := httptest.NewRecorder()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "This is a silly demo\n", string(body))
}
