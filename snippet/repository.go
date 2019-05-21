package snippet

import (
	"errors"
)

var (
	ErrNotFound = errors.New("snippet not found")
)

type Repository interface {
	Create(Snippet) error
	Update(Snippet) error
	Get(int) (Snippet, error)
	FindWithFilters(Filters) ([]Snippet, error)
	Del(int) error
}
