package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	w.WriteHeader(http.StatusCreated)

	fmt.Fprintln(w, "Title:", req.Title)
	fmt.Fprintln(w, "URL:", req.URL)
}
