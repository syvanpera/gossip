package snippet

import (
	"errors"
	"fmt"

	"github.com/syvanpera/gossip/ui"
)

const (
	COMMAND  = "COMMAND"
	CODE     = "CODE"
	BOOKMARK = "BOOKMARK"
	SNIP     = "SNIP"
)

var ErrNotExecutable = errors.New("not executable")
var ErrExecCanceled = errors.New("execution canceled")
var ErrEditCanceled = errors.New("edit canceled")

// SnippetData contains the data for a snippet
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

func New(data SnippetData) Snippet {
	switch data.Type {
	case COMMAND:
		return &Command{data: data}

	case CODE:
		return &Code{data: data}

	case BOOKMARK:
		return &Bookmark{data: data}
	}

	return nil
}

func Edit(label, s string) (string, error) {
	if s = ui.Prompt(label, s); s == "" {
		return "", ErrEditCanceled
	}

	return s, nil
}

func Editor(s string) (string, error) {
	if s = ui.Editor(s); s == "" {
		return "", ErrEditCanceled
	}

	return s, nil
}
