package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	url, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(url)

	p := make(map[string]string)
	p["message"] = "pong"
	json.NewEncoder(w).Encode(p)
}
