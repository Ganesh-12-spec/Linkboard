package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ganesh-12-spec/Linkboard/internal/database"
	"github.com/Ganesh-12-spec/Linkboard/internal/handler"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	pool, err := database.NewPool(ctx, connString)
	cancel()

	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	bookmarkHandler := handler.NewBookmarkHandler(pool)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /bookmarks", bookmarkHandler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", bookmarkHandler.CreateBookmark)

	// Serve the frontend from the same origin as the API.
	mux.Handle("/", http.FileServer(http.Dir("frontend")))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Linkboard server starting on :8080")

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
