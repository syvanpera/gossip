package snippet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"
)

const (
	COMMAND  = "COMMAND"
	CODE     = "CODE"
	BOOKMARK = "BOOKMARK"
	SNIP     = "SNIP"
)

var ErrNotExecutable = errors.New("Not executable")
var ErrExecCanceled = errors.New("Execution canceled")

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
	Edit(content, description string)
	String() string
}

type Snip struct {
	data SnippetData
}

func (*Snip) Type() SnippetType                { return SNIP }
func (s *Snip) Data() *SnippetData             { return &s.data }
func (*Snip) Execute() error                   { return ErrNotExecutable }
func (*Snip) Edit(content, description string) {}
func (s *Snip) String() string {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder
	width, _ := util.GetTerminalSize()
	description := runewidth.Truncate(s.data.Description, width-10, au.Gray(8, "...").String())
	border := strings.Repeat("─", width)
	borderVert := au.Gray(8, "│")

	fmt.Fprintln(&output, au.Gray(8, util.ReplaceRuneAtIndex(border, '┬', 8)))
	fmt.Fprintf(&output, "%s%s %s\n",
		au.Cyan(util.CenterStr(fmt.Sprintf("#%d", s.data.ID), 8)),
		borderVert,
		au.Yellow(description))
	fmt.Fprintln(&output, au.Gray(8, util.ReplaceRuneAtIndex(border, '┼', 8)))

	for i, s := range strings.Split(s.data.Content, "\n") {
		fmt.Fprintf(&output, "%s", au.Gray(8, util.CenterStr(strconv.Itoa(i+1), 8)))
		fmt.Fprintf(&output, "%s %s\n", borderVert, s)
	}
	fmt.Fprintln(&output, au.Gray(8, util.ReplaceRuneAtIndex(border, '┴', 8)))

	return output.String()
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

	default:
		return &Snip{data: sd}
	}
}
