package main

import (
	"fmt"
	"net/http"

	"github.com/Ganesh-12-spec/Linkboard/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	// Request flow:
	//
	// Client
	//   ↓
	// Router
	//   ↓
	// Handler
	//   ↓
	// Service
	//   ↓
	// Repository
	//   ↓
	// Storage
	//
	// The response travels back through the same layers.

	mux.HandleFunc("GET /bookmarks", handler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", handler.CreateBookmark)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
