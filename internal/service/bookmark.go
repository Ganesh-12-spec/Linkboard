package service

import "github.com/Ganesh-12-spec/Linkboard/internal/repository"

type Bookmark struct {
	ID    int
	Title string
	URL   string
}

type BookmarkService struct {
	repo repository.BookmarkRepository
}

func NewBookmarkService(repo repository.BookmarkRepository) BookmarkService {
	return BookmarkService{
		repo: repo,
	}
}

func (s BookmarkService) CreateBookmark(title string, url string) (Bookmark, error) {
	bookmark, err := s.repo.CreateBookmark(title, url)
	if err != nil {
		return Bookmark{}, err
	}

	return Bookmark{
		ID:    bookmark.ID,
		Title: bookmark.Title,
		URL:   bookmark.URL,
	}, nil
}