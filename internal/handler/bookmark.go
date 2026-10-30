package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
	"github.com/Ganesh-12-spec/Linkboard/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookmarkResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type CreateBookmarkRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type BookmarkHandler struct {
	service service.BookmarkService
}

func NewBookmarkHandler(db *pgxpool.Pool) *BookmarkHandler {
	repo := repository.NewBookmarkRepository(db)

	return &BookmarkHandler{
		service: service.NewBookmarkService(repo),
	}
}

func (h *BookmarkHandler) ListBookmarks(
	w http.ResponseWriter,
	r *http.Request,
) {
	bookmarks, err := h.service.ListBookmarks(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get bookmarks",
			http.StatusInternalServerError,
		)
		return
	}

	response := make([]BookmarkResponse, 0, len(bookmarks))

	for _, bookmark := range bookmarks {
		response = append(response, BookmarkResponse{
			ID:    bookmark.ID,
			Title: bookmark.Title,
			URL:   bookmark.URL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *BookmarkHandler) CreateBookmark(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateBookmarkRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.URL = strings.TrimSpace(req.URL)

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 255 {
		http.Error(w, "title must be 255 characters or fewer", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil ||
		parsedURL.Host == "" ||
		(parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "invalid HTTP or HTTPS URL", http.StatusBadRequest)
		return
	}

	bookmark, err := h.service.CreateBookmark(
		r.Context(),
		req.Title,
		req.URL,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create bookmark",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(BookmarkResponse{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	})
}
func (h BookmarkHandler) GetBookmark(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Step 1: Extract ID from URL and convert it to int64
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid bookmark ID", http.StatusBadRequest)
		return
	}

	// Step 2: Call Service layer
	bookmark, err := h.service.GetBookmark(r.Context(), id)
	if err != nil {
		// Step 3: Check if bookmark does not exist
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Bookmark not found", http.StatusNotFound)
			return
		}

		// Step 4: Handle unexpected errors
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Step 5: Send successful JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(bookmark); err != nil {
		return
	}
}
