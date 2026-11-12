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

// C19: List bookmarks with pagination
func (h *BookmarkHandler) ListBookmarks(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Read query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default pagination values
	limit := int64(10)
	offset := int64(0)

	// Validate limit
	if limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 64)

		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			http.Error(
				w,
				"limit must be between 1 and 100",
				http.StatusBadRequest,
			)
			return
		}

		limit = parsedLimit
	}

	// Validate offset
	if offsetStr != "" {
		parsedOffset, err := strconv.ParseInt(offsetStr, 10, 64)

		if err != nil || parsedOffset < 0 {
			http.Error(
				w,
				"offset must be zero or greater",
				http.StatusBadRequest,
			)
			return
		}

		offset = parsedOffset
	}

	// Call Service layer
	bookmarks, err := h.service.ListBookmarks(
		r.Context(),
		limit,
		offset,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get bookmarks",
			http.StatusInternalServerError,
		)
		return
	}

	// Convert service objects into response objects
	response := make([]BookmarkResponse, 0, len(bookmarks))

	for _, bookmark := range bookmarks {
		response = append(response, BookmarkResponse{
			ID:    bookmark.ID,
			Title: bookmark.Title,
			URL:   bookmark.URL,
		})
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

// C17: Create bookmark
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

// C18: Get one bookmark by ID
func (h *BookmarkHandler) GetBookmark(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid bookmark ID", http.StatusBadRequest)
		return
	}

	bookmark, err := h.service.GetBookmark(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Bookmark not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(BookmarkResponse{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	})
}
func (h *BookmarkHandler) UpdateBookmark(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid bookmark ID", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateBookmarkRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// 3. Validate input
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

	// 4. Call service layer
	bookmark, err := h.service.UpdateBookmark(
		r.Context(),
		id,
		req.Title,
		req.URL,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Bookmark not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. Send updated bookmark
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(BookmarkResponse{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	})
}
