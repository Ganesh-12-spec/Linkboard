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

// C17: Create bookmark
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

// C19: List bookmarks with pagination
func (s BookmarkService) ListBookmarks(
	ctx context.Context,
	limit, offset int64,
) ([]Bookmark, error) {
	bookmarks, err := s.repo.ListBookmarks(ctx, limit, offset)
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

// C18: Get bookmark by ID
func (s BookmarkService) GetBookmark(
	ctx context.Context,
	id int64,
) (Bookmark, error) {
	bookmark, err := s.repo.GetBookmark(ctx, id)
	if err != nil {
		return Bookmark{}, err
	}

	return Bookmark{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	}, nil
}
func(s BookmarkService) UpdateBookmark(
	ctx context.Context,
	id int64,
	title, url string,

)(Bookmark, error){
	bookmark, err := s.repo.UpdateBookmark(ctx,id, title, url)
	if err != nil {
		return Bookmark{}, err

	}
	return Bookmark{

		Title: bookmark.Title,
		 URL: bookmark.URL,
	}, nil
}
func (s BookmarkService) DeleteBookmark(
    ctx context.Context,
    id int64,
) error {
    err := s.repo.DeleteBookmark(ctx, id)
    if err != nil {
        return err
    }

    return nil
}