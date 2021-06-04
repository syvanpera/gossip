package snippet

import (
	"errors"
	"fmt"
)

const (
	COMMAND  = "Commands"
	CODE     = "Code snippets"
	BOOKMARK = "Bookmarks"
	SNIP     = "SNIP"
)

var ErrNotExecutable = errors.New("not executable")
var ErrExecCanceled = errors.New("execution canceled")
var ErrEditCanceled = errors.New("edit canceled")

// SnippetData contains the data for one snippet
type SnippetData struct {
	ID          int64
	Content     string
	Description string
	Tags        string
	Type        SnippetType
	Language    string
}

type SnippetType string

type Snippet interface {
	Type() SnippetType
	Data() *SnippetData
	Execute() error
	Edit() error
	String() string
}

// Filters are used to filter the snippets list
type Filters struct {
	Type     SnippetType
	Language string
	Tags     string
}

func (f Filters) String() string {
	return fmt.Sprintf("{Type: \"%s\", Language: \"%s\", Tags: \"%s\"}", f.Type, f.Language, f.Tags)
}

func New(sd SnippetData) Snippet {
	switch sd.Type {
	case COMMAND:
		return &Command{data: sd}

	case CODE:
		return &Code{data: sd}

	case BOOKMARK:
		return &Bookmark{data: sd}
	}

	return nil
}
