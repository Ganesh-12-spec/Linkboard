package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookmarkRepository struct {
	db *pgxpool.Pool
}

type Bookmark struct {
	ID    int64
	Title string
	URL   string
}

func NewBookmarkRepository(db *pgxpool.Pool) BookmarkRepository {
	return BookmarkRepository{db: db}
}

func (r BookmarkRepository) CreateBookmark(
	ctx context.Context,
	title string,
	url string,
) (Bookmark, error) {
	var bookmark Bookmark

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO bookmarks (title, url)
		 VALUES ($1, $2)
		 RETURNING id, title, url`,
		title,
		url,
	).Scan(&bookmark.ID, &bookmark.Title, &bookmark.URL)

	if err != nil {
		return Bookmark{}, err
	}

	return bookmark, nil
}

func (r BookmarkRepository) ListBookmarks(
	ctx context.Context,
) ([]Bookmark, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, title, url
		 FROM bookmarks
		 ORDER BY created_at DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookmarks := make([]Bookmark, 0)

	for rows.Next() {
		var bookmark Bookmark

		if err := rows.Scan(
			&bookmark.ID,
			&bookmark.Title,
			&bookmark.URL,
		); err != nil {
			return nil, err
		}

		bookmarks = append(bookmarks, bookmark)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookmarks, nil
}
