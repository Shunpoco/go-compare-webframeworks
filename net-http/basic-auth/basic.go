package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var db = make(map[string]string)

func setupRouter() *http.ServeMux {
	r := http.NewServeMux()

	// ping test
	r.HandleFunc("/ping", pingHandler)

	// Get user value
	r.HandleFunc("/user/", userHandler)

	authorized := make(map[string]string)

	authorized["foo"] = "bar"
	authorized["manu"] = "123"

	adminHandler := makeAdminHandler(authorized)

	// Basic Auth
	r.HandleFunc("/admin", adminHandler)

	return r
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := strings.TrimPrefix(r.URL.Path, "/user/")

	log.Println(user)

	w.Header().Set("Content-Type", "application/json")
	p := make(map[string]string)

	value, ok := db[user]
	p["user"] = user
	if ok {
		p["value"] = value
	} else {
		p["status"] = "no value"
	}
	json.NewEncoder(w).Encode(p)
	w.WriteHeader(http.StatusOK)
}

func makeAdminHandler(authorized map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if p, ok := authorized[user]; ok && p == pass {
			db[user] = pass
			w.WriteHeader(http.StatusOK)
			var result = map[string]string{"status": "ok"}
			json.NewEncoder(w).Encode(result)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func main() {
	serveMux := setupRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
