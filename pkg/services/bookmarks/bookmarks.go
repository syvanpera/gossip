package bookmarks

import (
	"fmt"
	"strings"
)

// Repo defines the DB level interaction of bookmarks
type Repo interface {
	Get(id int) (Bookmark, error)
	GetAll(filters BookmarkFilters) ([]Bookmark, error)
	Create(bcu BookmarkCreateUpdate) (int, error)
	Update(bcu BookmarkCreateUpdate, id int) error
	Delete(id int) error
	InitDB()
}

// Service defines the service level contract that other services
// outside this package can use to interact with Bookmark resources
type Service interface {
	Get(id int) (Bookmark, error)
	Find(tags string) ([]Bookmark, error)
	Create(bcu BookmarkCreateUpdate) (Bookmark, error)
	Update(bcu BookmarkCreateUpdate, id int) (Bookmark, error)
	Delete(id int) error
}

type bookmark struct {
	repo Repo
}

func New(repo Repo) Service {
	return &bookmark{repo}
}

func (s *bookmark) Get(id int) (Bookmark, error) {
	return s.repo.Get(id)
}

func (s *bookmark) Find(tags string) ([]Bookmark, error) {
	var tagFilters []string

	if tags != "" {
		tags := strings.Split(tags, ",")
		for _, tag := range tags {
			tagFilters = append(tagFilters, fmt.Sprintf("',' || tags || ',' like '%%,%s,%%'", tag))
		}
	}

	filters := BookmarkFilters{
		Tags: tagFilters,
	}

	return s.repo.GetAll(filters)
}

func (s *bookmark) Create(bcu BookmarkCreateUpdate) (Bookmark, error) {
	id, err := s.repo.Create(bcu)
	if err != nil {
		return Bookmark{}, err
	}
	return s.repo.Get(id)
}

func (s *bookmark) Update(bcu BookmarkCreateUpdate, id int) (Bookmark, error) {
	if err := s.repo.Update(bcu, id); err != nil {
		return Bookmark{}, err
	}
	return s.Get(id)
}

func (s *bookmark) Delete(id int) error {
	return s.repo.Delete(id)
}
