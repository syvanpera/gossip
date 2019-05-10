package snippet

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
)

// Snippet contains the data for one snippet
type Snippet struct {
	ID          int      `json:"id"`
	Snippet     string   `json:"snippet"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Language    string   `json:"language,omitempty"`
}

func (s Snippet) String() string {
	snippet := "Snippet:"
	switch s.Type {
	case "code":
		snippet = "Code:"
	case "cmd":
		snippet = "Command:"
	case "url":
		snippet = "URL:"
	}
	tags := strings.Join(s.Tags, " ")

	return fmt.Sprintf("%12s %d\n%12s %s\n%12s %s\n%12s %s\n%12s %s",
		aurora.Gray(10, "ID:"), s.ID,
		aurora.Yellow("Description:"), s.Description,
		aurora.Green(snippet), s.Snippet,
		aurora.Cyan("Tags:"), tags,
		aurora.Blue("Language:"), s.Language)
}
