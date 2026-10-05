package main

import (
	"fmt"
	"net/http"

	"github.com/Ganesh-12-spec/Linkboard/internal/handler"
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

	mux.HandleFunc("GET /bookmarks", handler.GetBookmark)

	mux.HandleFunc("POST /bookmarks", handler.CreateBookmark)

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
