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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := database.NewPool(ctx, connString)
	cancel()

	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /bookmarks", handler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", handler.CreateBookmark)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Linkboard server starting on :8080")

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("server error: %v", err)
	}
}
