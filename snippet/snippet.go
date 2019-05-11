package snippet

import (
	"fmt"
	"strings"
)

// Snippet contains the data for one snippet
type Snippet struct {
	ID          int      `json:"id"`
	Snippet     string   `json:"snippet"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Language    string   `json:"language,omitempty"`
}

// Filters are used to filter the snippets list
type Filters struct {
	Type     string
	Language string
	Tags     []string
}

func (s Snippet) String() string {
	switch s.Type {
	case "code", "cmd":
		return renderCode(s)
	default:
		return render(s)
	}
}

func (f Filters) String() string {
	return fmt.Sprintf("{Type: \"%s\", Language: \"%s\", Tags: \"%s\"}", f.Type, f.Language, strings.Join(f.Tags, ","))
}
