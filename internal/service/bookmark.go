package service

type Bookmark struct {
	ID    int
	Title string
	URL   string
}

type BookmarkService struct{}

func (s BookmarkService) CreateBookmark(title string, url string) (Bookmark, error) {
	bookmark := Bookmark{
		ID:    1,
		Title: title,
		URL:   url,
	}

	return bookmark, nil
}
