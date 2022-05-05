package bookmarks

import (
	"errors"
)

var (
	ErrBookmarkNotFound = errors.New("requested bookmark could not be found")
	ErrBookmarkQuery    = errors.New("requested bookmarks could not be retrieved based on the given criteria")
	ErrBookmarkCreate   = errors.New("bookmark could not be created")
	ErrBookmarkUpdate   = errors.New("bookmark could not be updated")
	ErrBookmarkDelete   = errors.New("bookmark could not be deleted")
)
