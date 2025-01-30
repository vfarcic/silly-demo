//go:build integration

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func TestVideoPut(t *testing.T) {
	t.Run("should put a video into the database", func(t *testing.T) {
		// Test case 1: Post a video
		rand.Seed(time.Now().UnixNano())
		expectedID := strconv.Itoa(rand.Intn(100000))
		expectedTitle := make([]byte, 100)
		for i := range expectedTitle {
			expectedTitle[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		url := fmt.Sprintf("http://silly-demo.127.0.0.1.nip.io/video?id=%s&title=%s", expectedID, expectedTitle)
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer res.Body.Close()
		w := httptest.NewRecorder()
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		// Test case 2: Get videos
		url = "http://silly-demo.127.0.0.1.nip.io/videos"
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err = client.Do(req)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer res.Body.Close()
		w = httptest.NewRecorder()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf(err.Error())
		}
		var videos []Video
		err = json.Unmarshal(body, &videos)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, videos)
		found := false
		for _, video := range videos {
			if video.ID == expectedID && video.Title == string(expectedTitle) {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("Expected video ID %s and title %s not found in the response", expectedID, expectedTitle))
	})
}
