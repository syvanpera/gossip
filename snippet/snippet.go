package snippet

import (
	"database/sql"
	"fmt"
)

// Snippet contains the data for one snippet
type Snippet struct {
	ID          int            `json:"id"`
	Snippet     string         `json:"snippet"`
	Description sql.NullString `json:"description,omitempty"`
	Tags        []string       `json:"tags"`
	Type        string         `json:"type"`
	Language    sql.NullString `json:"language,omitempty"`
}

func (s Snippet) String() string {
	return fmt.Sprintf("%d: %s", s.ID, s.Snippet)
}
