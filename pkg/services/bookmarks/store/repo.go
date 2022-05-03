package store

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/syvanpera/gossip/pkg/services/bookmarks"
)

const (
	selectBookmark      = `SELECT id, url, description, tags, flags, created_at, updated_at FROM bookmarks WHERE id=$1`
	selectManyBookmarks = `SELECT id, url, description, tags, flags, created_at, updated_at FROM bookmarks`
	insertBookmark      = `INSERT INTO bookmarks (url, description, tags, flags) VALUES ($1, $2, $3, $4) RETURNING id`
	updateBookmark      = `UPDATE bookmarks SET url = $1, description = $2, tags = $3, flags = $4, updated_at = datetime('now') WHERE id = $5`
	deleteBookmark      = `DELETE FROM bookmarks WHERE id = $1`
)

type bookmarkRepo struct {
	DB *sql.DB
}

func New(conn *sql.DB) bookmarks.Repo {
	return &bookmarkRepo{conn}
}

func (r *bookmarkRepo) Get(id int) (bookmarks.Bookmark, error) {
	var bm bookmarks.Bookmark

	err := r.DB.QueryRow(selectBookmark, id).
		Scan(&bm.ID, &bm.URL, &bm.Description, &bm.Tags, &bm.Flags, &bm.CreatedAt, &bm.UpdatedAt)
	if err != nil {
		log.Err(err).Msg("Getting bookmark failed")
		return bm, bookmarks.ErrBookmarkNotFound
	}

	return bm, nil
}

func (r *bookmarkRepo) GetAll() ([]bookmarks.Bookmark, error) {
	bl := make([]bookmarks.Bookmark, 0)

	rows, err := r.DB.Query(selectManyBookmarks)
	if err != nil {
		log.Err(err).Msg("Getting bookmarks failed")
		return bl, bookmarks.ErrBookmarkQuery
	}
	defer rows.Close()

	for rows.Next() {
		var bm bookmarks.Bookmark
		if err := rows.Scan(&bm.ID, &bm.URL, &bm.Description, &bm.Tags, &bm.Flags, &bm.CreatedAt, &bm.UpdatedAt); err != nil {
			log.Err(err).Msg("Getting bookmarks failed")
			return bl, bookmarks.ErrBookmarkQuery
		}

		bl = append(bl, bm)
	}

	return bl, nil
}

func (r *bookmarkRepo) Create(bcu bookmarks.BookmarkCreateUpdate) (int, error) {
	var id int
	if err := r.DB.QueryRow(insertBookmark, bcu.URL, bcu.Description, bcu.Tags, bcu.Flags).Scan(&id); err != nil {
		log.Err(err).Msg("Creating bookmark failed")
		return -1, bookmarks.ErrBookmarkCreate
	}

	return id, nil
}

func (r *bookmarkRepo) Update(bcu bookmarks.BookmarkCreateUpdate, id int) error {
	_, err := r.DB.Exec(updateBookmark, bcu.URL, bcu.Description, bcu.Tags, bcu.Flags, id)
	if err != nil {
		log.Err(err).Int("ID", id).Msg("Updating bookmark failed")
		return bookmarks.ErrBookmarkUpdate
	}
	return nil
}

func (r *bookmarkRepo) Delete(id int) error {
	_, err := r.DB.Exec(deleteBookmark, id)
	if err != nil {
		log.Err(err).Msg("Deleting bookmark failed")
		return bookmarks.ErrBookmarkDelete
	}

	return nil
}

func (r *bookmarkRepo) InitDB() {
	schema := `
CREATE TABLE IF NOT EXISTS bookmarks(
    id INTEGER NOT NULL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    flags INTEGER NOT NULL DEFAULT 0,

    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
)`
	r.DB.Exec(schema)
}
