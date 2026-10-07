package repository

type BookmarkRepository struct{}

type Bookmark struct {
	ID    int
	Title string
	URL   string
}

func (r BookmarkRepository) CreateBookmark(title string, url string) (Bookmark, error) {
	bookmark := Bookmark{
		ID:    1,
		Title: title,
		URL:   url,
	}

	return bookmark, nil
}