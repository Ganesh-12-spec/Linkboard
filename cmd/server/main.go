package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Linkboard")
	})

	mux.HandleFunc("GET /bookmarks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bookmarks")
	})

	mux.HandleFunc("POST /bookmarks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create bookmark")
	})

	mux.HandleFunc("GET /bookmarks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintln(w, "Bookmark ID:", id)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
