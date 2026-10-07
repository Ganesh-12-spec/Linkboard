package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
	"github.com/Ganesh-12-spec/Linkboard/internal/service"
)

type BookmarkResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type CreateBookmarkRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func GetBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark := BookmarkResponse{
		ID:    1,
		Title: "Learn Go",
		URL:   "https://go.dev",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(bookmark)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req CreateBookmarkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	logRequestContext(ctx)

	bookmarkRepository := repository.BookmarkRepository{}
	bookmarkService := service.NewBookmarkService(bookmarkRepository)

	bookmark, err := bookmarkService.CreateBookmark(req.Title, req.URL)
	if err != nil {
		http.Error(w, "failed to create bookmark", http.StatusInternalServerError)
		return
	}

	response := BookmarkResponse{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func logRequestContext(ctx context.Context) {
	fmt.Println("Request context received")
}
