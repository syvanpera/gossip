package snippet

import (
	"errors"
	"fmt"
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
	String() string
}

type Code struct {
	data SnippetData
}

func (c *Code) Type() SnippetType  { return CODE }
func (c *Code) Data() *SnippetData { return &c.data }
func (c *Code) Execute() error     { return ErrNotExecutable }
func (c *Code) String() string {
	return renderCode(c.data)
}

type Command struct {
	data SnippetData
}

func (cmd *Command) Type() SnippetType  { return COMMAND }
func (cmd *Command) Data() *SnippetData { return &cmd.data }

func (cmd *Command) Execute() error {
	fmt.Println("Okay, executing command...")
	return nil
}

func (cmd *Command) String() string {
	return render(cmd.data)
}

type Snip struct {
	data SnippetData
}

func (s *Snip) Type() SnippetType  { return SNIP }
func (s *Snip) Data() *SnippetData { return &s.data }
func (s *Snip) Execute() error     { return ErrNotExecutable }
func (s *Snip) String() string     { return render(s.data) }

// Filters are used to filter the snippets list
type Filters struct {
	Type     string
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

	default:
		return &Snip{data: sd}
	}
}
