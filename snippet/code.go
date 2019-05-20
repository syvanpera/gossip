package snippet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"
)

var style = styles.Register(chroma.MustNewStyle("gruvbox", chroma.StyleEntries{
	chroma.Comment:                  "#7c6f64",
	chroma.CommentHashbang:          "#7c6f64",
	chroma.CommentMultiline:         "#7c6f64",
	chroma.CommentPreproc:           "#fb4933",
	chroma.CommentSingle:            "#7c6f64",
	chroma.CommentSpecial:           "#7c6f64",
	chroma.Generic:                  "#fdf4c1",
	chroma.GenericDeleted:           "#8b080b",
	chroma.GenericEmph:              "#fdf4c1 underline",
	chroma.GenericError:             "#fdf4c1",
	chroma.GenericHeading:           "#fdf4c1 bold",
	chroma.GenericInserted:          "#fdf4c1 bold",
	chroma.GenericOutput:            "#44475a",
	chroma.GenericPrompt:            "#fdf4c1",
	chroma.GenericStrong:            "#fdf4c1",
	chroma.GenericSubheading:        "#fdf4c1 bold",
	chroma.GenericTraceback:         "#fdf4c1",
	chroma.GenericUnderline:         "underline",
	chroma.Error:                    "#fdf4c1",
	chroma.Keyword:                  "#fb4933",
	chroma.KeywordConstant:          "#fb4933",
	chroma.KeywordDeclaration:       "#fb4933 italic",
	chroma.KeywordNamespace:         "#fb4933",
	chroma.KeywordPseudo:            "#fb4933",
	chroma.KeywordReserved:          "#fb4933",
	chroma.KeywordType:              "#d3869b",
	chroma.Literal:                  "#fdf4c1",
	chroma.LiteralDate:              "#fdf4c1",
	chroma.Name:                     "#fdf4c1",
	chroma.NameAttribute:            "#fabd2f",
	chroma.NameBuiltin:              "#fb4933 italic",
	chroma.NameBuiltinPseudo:        "#fdf4c1",
	chroma.NameClass:                "#fabd2f",
	chroma.NameConstant:             "#fdf4c1",
	chroma.NameDecorator:            "#fdf4c1",
	chroma.NameEntity:               "#fdf4c1",
	chroma.NameException:            "#fdf4c1",
	chroma.NameFunction:             "#fabd2f",
	chroma.NameLabel:                "#fb4933 italic",
	chroma.NameNamespace:            "#fdf4c1",
	chroma.NameOther:                "#fdf4c1",
	chroma.NameTag:                  "#fb4933",
	chroma.NameVariable:             "#fb4933 italic",
	chroma.NameVariableClass:        "#fb4933 italic",
	chroma.NameVariableGlobal:       "#fb4933 italic",
	chroma.NameVariableInstance:     "#fb4933 italic",
	chroma.LiteralNumber:            "#bd93f9",
	chroma.LiteralNumberBin:         "#bd93f9",
	chroma.LiteralNumberFloat:       "#bd93f9",
	chroma.LiteralNumberHex:         "#bd93f9",
	chroma.LiteralNumberInteger:     "#bd93f9",
	chroma.LiteralNumberIntegerLong: "#bd93f9",
	chroma.LiteralNumberOct:         "#bd93f9",
	chroma.Operator:                 "#fb4933",
	chroma.OperatorWord:             "#fb4933",
	chroma.Other:                    "#fdf4c1",
	chroma.Punctuation:              "#fdf4c1",
	chroma.LiteralString:            "#b8bb26",
	chroma.LiteralStringBacktick:    "#b8bb26",
	chroma.LiteralStringChar:        "#b8bb26",
	chroma.LiteralStringDoc:         "#b8bb26",
	chroma.LiteralStringDouble:      "#b8bb26",
	chroma.LiteralStringEscape:      "#b8bb26",
	chroma.LiteralStringHeredoc:     "#b8bb26",
	chroma.LiteralStringInterpol:    "#b8bb26",
	chroma.LiteralStringOther:       "#b8bb26",
	chroma.LiteralStringRegex:       "#b8bb26",
	chroma.LiteralStringSingle:      "#b8bb26",
	chroma.LiteralStringSymbol:      "#b8bb26",
	chroma.Text:                     "#fdf4c1",
	chroma.TextWhitespace:           "#fdf4c1",
}))

type Code struct {
	data SnippetData
}

func (c *Code) Type() SnippetType  { return CODE }
func (c *Code) Data() *SnippetData { return &c.data }

func (c *Code) Execute() error {
	return ErrNotExecutable
}

func (c *Code) Edit() error {
	s, err := Editor(c.Data().Content)
	if err != nil {
		return err
	}
	c.Data().Content = s

	s, err = Edit("Description", c.Data().Description)
	if err != nil {
		return err
	}
	c.Data().Description = s

	return nil
}

func (c *Code) String() string {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder
	width, _ := util.GetTerminalSize()
	description := runewidth.Truncate(c.data.Description, width-10, au.Gray(8, "...").String())
	border := strings.Repeat("─", width)
	borderVert := au.Gray(8, "│")

	fmt.Fprintf(&output, "\n%s\n", au.Gray(8, util.ReplaceRuneAtIndex(border, '┬', 8)))
	fmt.Fprintf(&output, "%s%s %s\n",
		au.Cyan(util.CenterStr(fmt.Sprintf("#%d", c.data.ID), 8)),
		borderVert,
		au.Yellow(description))
	fmt.Fprintln(&output, au.Gray(8, util.ReplaceRuneAtIndex(border, '┼', 8)))

	content := c.data.Content
	if colors {
		var sb strings.Builder
		quick.Highlight(&sb, c.data.Content, c.data.Language, "terminal16m", "gruvbox")
		content = sb.String()
	}

	for i, s := range strings.Split(content, "\n") {
		fmt.Fprintf(&output, "%s", au.Gray(8, util.CenterStr(strconv.Itoa(i+1), 8)))
		fmt.Fprintf(&output, "%s %s\n", borderVert, s)
	}
	fmt.Fprintf(&output, "%s", au.Gray(8, util.ReplaceRuneAtIndex(border, '┴', 8)))

	return output.String()
}

func NewCode(content, description string, tags string, language string) Code {
	cmd := Code{
		data: SnippetData{
			Content:     content,
			Description: description,
			Tags:        strings.ToLower(tags),
			Language:    strings.ToLower(language),
			Type:        CODE,
		},
	}

	return cmd
}
