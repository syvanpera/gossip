package snippet

type Repository interface {
	Create(Snippet) error
	Update(Snippet) error
	Get(int) (Snippet, error)
	FindWithFilters(Filters) ([]Snippet, error)
	Delete(int) error
}
