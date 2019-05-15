package snippet

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/browser"
)

const (
	COMMAND  = "COMMAND"
	CODE     = "CODE"
	BOOKMARK = "BOOKMARK"
	SNIP     = "SNIP"
)

var ErrNotExecutable = errors.New("Not executable")

// SnippetData contains the data for one snippet
type SnippetData struct {
	ID          int         `json:"id"`
	Snippet     string      `json:"snippet"`
	Description string      `json:"description"`
	Tags        []string    `json:"tags"`
	Type        SnippetType `json:"type"`
	Language    string      `json:"language,omitempty"`
}

type SnippetType string

type Snippet interface {
	Type() SnippetType
	Execute() error
	String() string
}

type Code struct {
	Data SnippetData
}

func (c *Code) Type() SnippetType { return CODE }
func (c *Code) Execute() error    { return ErrNotExecutable }
func (c *Code) String() string {
	return renderCode(c.Data)
}

type Command struct {
	Data SnippetData
}

func (cmd *Command) Type() SnippetType { return COMMAND }

func (cmd *Command) Execute() error {
	fmt.Println("Okay, executing command...")
	return nil
}

func (cmd *Command) String() string {
	return render(cmd.Data)
}

type Bookmark struct {
	Data SnippetData
}

func (b *Bookmark) Type() SnippetType { return BOOKMARK }

func (b *Bookmark) Execute() error {
	fmt.Println("Okay, opening link in default browser...")
	url := b.Data.Snippet
	if matched, _ := regexp.MatchString(`^http(s)?://*`, url); !matched {
		url = "http://" + url
	}
	browser.OpenURL(url)

	return nil
}

func (b *Bookmark) String() string {
	return renderBookmark(b.Data)
}

type Snip struct {
	Data SnippetData
}

func (s *Snip) Type() SnippetType { return SNIP }
func (s *Snip) Execute() error    { return ErrNotExecutable }
func (s *Snip) String() string {
	return render(s.Data)
}

// Filters are used to filter the snippets list
type Filters struct {
	Type     string
	Language string
	Tags     []string
}

func (f Filters) String() string {
	return fmt.Sprintf("{Type: \"%s\", Language: \"%s\", Tags: \"%s\"}", f.Type, f.Language, strings.Join(f.Tags, ","))
}

func New(sd SnippetData) Snippet {
	switch sd.Type {
	case COMMAND:
		return &Command{Data: sd}

	case CODE:
		return &Code{Data: sd}

	case BOOKMARK:
		return &Bookmark{Data: sd}

	default:
		return &Snip{Data: sd}
	}
}
