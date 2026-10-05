package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintln(w, "Title:", req.Title)
	fmt.Fprintln(w, "URL:", req.URL)
}

func logRequestContext(ctx context.Context) {
	fmt.Println("Request context received")
}