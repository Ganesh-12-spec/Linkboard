package service

import (
	"context"

	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
)

type Bookmark struct {
	ID    int64
	Title string
	URL   string
}

type BookmarkService struct {
	repo repository.BookmarkRepository
}

func NewBookmarkService(
	repo repository.BookmarkRepository,
) BookmarkService {
	return BookmarkService{repo: repo}
}

func (s BookmarkService) CreateBookmark(
	ctx context.Context,
	title string,
	url string,
) (Bookmark, error) {
	bookmark, err := s.repo.CreateBookmark(ctx, title, url)
	if err != nil {
		return Bookmark{}, err
	}

	return Bookmark{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	}, nil
}

func (s BookmarkService) ListBookmarks(
	ctx context.Context,
) ([]Bookmark, error) {
	bookmarks, err := s.repo.ListBookmarks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Bookmark, 0, len(bookmarks))

	for _, bookmark := range bookmarks {
		result = append(result, Bookmark{
			ID:    bookmark.ID,
			Title: bookmark.Title,
			URL:   bookmark.URL,
		})
	}

	return result, nil
}
