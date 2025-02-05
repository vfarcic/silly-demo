//go:build integration

package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibonacci(t *testing.T) {
	// Test case 1: Valid number
	url := "http://silly-demo.127.0.0.1.nip.io/fibonacci?number=5"
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, `{"number":5,"fibonacci":5}`, string(body))

	// Test case 2: Invalid input
	url = "http://silly-demo.127.0.0.1.nip.io/fibonacci?number=abc"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	res, err = client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax"}`, string(body))

	// Test case 3: Missing number parameter
	url = "http://silly-demo.127.0.0.1.nip.io/fibonacci"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	res, err = client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, `{"error":"strconv.Atoi: parsing \"\": invalid syntax"}`, string(body))
}
