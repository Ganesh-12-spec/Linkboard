package main

import (
	"fmt"
	"net/http"
)

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from a struct handler")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Linkboard")
	})

	mux.HandleFunc("GET /bookmarks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Method:", r.Method)
		fmt.Fprintln(w, "Path:", r.URL.Path)
		fmt.Fprintln(w, "User-Agent:", r.Header.Get("User-Agent"))
	})

	mux.HandleFunc("POST /bookmarks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Method:", r.Method)
		fmt.Fprintln(w, "Content-Type:", r.Header.Get("Content-Type"))
		fmt.Fprintln(w, "Request received for creating a bookmark")
	})

	mux.HandleFunc("GET /bookmarks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		fmt.Fprintln(w, "Method:", r.Method)
		fmt.Fprintln(w, "Bookmark ID:", id)
	})

	mux.Handle("GET /hello", helloHandler{})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
