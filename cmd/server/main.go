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
	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
	"github.com/Ganesh-12-spec/Linkboard/internal/service"
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
	defer cancel()

	pool, err := database.NewPool(ctx, connString)
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	// Bookmark handler
	bookmarkHandler := handler.NewBookmarkHandler(pool)

	// User registration
	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /bookmarks", bookmarkHandler.ListBookmarks)
	mux.HandleFunc("GET /bookmarks/{id}", bookmarkHandler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", bookmarkHandler.CreateBookmark)
	mux.HandleFunc("PUT /bookmarks/{id}", bookmarkHandler.UpdateBookmark)
	mux.HandleFunc("DELETE /bookmarks/{id}", bookmarkHandler.DeleteBookmark)

	mux.HandleFunc("POST /auth/register", authHandler.Register)

	// Serve frontend
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
