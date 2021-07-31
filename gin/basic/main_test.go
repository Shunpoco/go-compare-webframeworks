package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AdminResponseBody struct {
	Status string `json:"status"`
}

type User struct {
	User   string `json:"user"`
	Value  string `json:"value"`
	Status string `json:"status"`
}

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestAdminRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonStr = []byte(`{"value": "bar"}`)
	req, _ := http.NewRequest("POST", "/admin", bytes.NewBuffer(jsonStr))
	req.SetBasicAuth("foo", "bar")
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	var arb AdminResponseBody
	json.Unmarshal(w.Body.Bytes(), &arb)

	assert.Equal(t, "ok", arb.Status)
}

func TestGetUser(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	var jsonStr = []byte(`{"value": "bar"}`)
	req, _ := http.NewRequest("POST", "/admin", bytes.NewBuffer(jsonStr))
	req.SetBasicAuth("foo", "bar")
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/foo", nil)
	router.ServeHTTP(w, req)

	var user User
	json.Unmarshal(w.Body.Bytes(), &user)

	assert.Equal(t, "bar", user.Value)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/manu", nil)
	router.ServeHTTP(w, req)

	var user2 User
	json.Unmarshal(w.Body.Bytes(), &user2)

	assert.Equal(t, "no value", user2.Status)
}
