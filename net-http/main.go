package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", ping)
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Fatal(server.ListenAndServe())
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := make(map[string]string)
	p["message"] = "pong"
	json.NewEncoder(w).Encode(p)
}
