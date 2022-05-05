package service

import (
	"fmt"
	"strings"

	"github.com/syvanpera/gossip/pkg/bookmarks"
	"github.com/syvanpera/gossip/pkg/bookmarks/store"
)

// Service defines the service level contract that other services
// outside this package can use to interact with Bookmark resources
type Service interface {
	Get(id int) (bookmarks.Bookmark, error)
	Find(tags string) ([]bookmarks.Bookmark, error)
	Create(bcu bookmarks.BookmarkCreateUpdate) (bookmarks.Bookmark, error)
	Update(bcu bookmarks.BookmarkCreateUpdate, id int) (bookmarks.Bookmark, error)
	Delete(id int) error
}

type bookmark struct {
	repo store.Repo
}

func New(repo store.Repo) Service {
	return &bookmark{repo}
}

func (s *bookmark) Get(id int) (bookmarks.Bookmark, error) {
	return s.repo.Get(id)
}

func (s *bookmark) Find(tags string) ([]bookmarks.Bookmark, error) {
	var tagFilters []string

	if tags != "" {
		tags := strings.Split(tags, ",")
		for _, tag := range tags {
			tagFilters = append(tagFilters, fmt.Sprintf("',' || tags || ',' like '%%,%s,%%'", tag))
		}
	}

	filters := bookmarks.BookmarkFilters{
		Tags: tagFilters,
	}

	return s.repo.GetAll(filters)
}

func (s *bookmark) Create(bcu bookmarks.BookmarkCreateUpdate) (bookmarks.Bookmark, error) {
	id, err := s.repo.Create(bcu)
	if err != nil {
		return bookmarks.Bookmark{}, err
	}
	return s.repo.Get(id)
}

func (s *bookmark) Update(bcu bookmarks.BookmarkCreateUpdate, id int) (bookmarks.Bookmark, error) {
	if err := s.repo.Update(bcu, id); err != nil {
		return bookmarks.Bookmark{}, err
	}
	return s.Get(id)
}

func (s *bookmark) Delete(id int) error {
	return s.repo.Delete(id)
}
