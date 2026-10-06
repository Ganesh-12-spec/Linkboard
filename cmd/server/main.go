package main

import (
	"fmt"
	"net/http"

	"github.com/Ganesh-12-spec/Linkboard/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	// Router decides WHERE the HTTP request should go.
	// The actual HTTP handling lives inside the handler package.
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
