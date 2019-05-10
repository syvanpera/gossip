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

func renderCode(s Snippet) string {
	colors := viper.GetBool("colorOutput") != viper.GetBool("color")
	au := aurora.NewAurora(colors)
	var output strings.Builder
	width, _ := GetTerminalSize()
	description := runewidth.Truncate(s.Description, width-10, au.Gray(8, "...").String())
	border := strings.Repeat("─", width)
	borderVert := au.Gray(8, "│")

	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┬', 8)))
	fmt.Fprintf(&output, "%s%s %s\n",
		au.Cyan(CenterStr(fmt.Sprintf("#%d", s.ID), 8)),
		borderVert,
		au.Yellow(description))
	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┼', 8)))

	snippet := s.Snippet
	if colors {
		var sb strings.Builder
		quick.Highlight(&sb, s.Snippet, s.Language, "terminal16m", "gruvbox")
		snippet = sb.String()
	}

	for i, s := range strings.Split(snippet, "\n") {
		fmt.Fprintf(&output, "%s", au.Gray(8, CenterStr(strconv.Itoa(i+1), 8)))
		fmt.Fprintf(&output, "%s %s\n", borderVert, s)
	}
	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┴', 8)))

	return output.String()
}

func render(s Snippet) string {
	colors := viper.GetBool("colorOutput") != viper.GetBool("color")
	au := aurora.NewAurora(colors)
	var output strings.Builder
	width, _ := GetTerminalSize()
	description := runewidth.Truncate(s.Description, width-10, au.Gray(8, "...").String())
	border := strings.Repeat("─", width)
	borderVert := au.Gray(8, "│")

	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┬', 8)))
	fmt.Fprintf(&output, "%s%s %s\n",
		au.Cyan(CenterStr(fmt.Sprintf("#%d", s.ID), 8)),
		borderVert,
		au.Yellow(description))
	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┼', 8)))

	for i, s := range strings.Split(s.Snippet, "\n") {
		fmt.Fprintf(&output, "%s", au.Gray(8, CenterStr(strconv.Itoa(i+1), 8)))
		fmt.Fprintf(&output, "%s %s\n", borderVert, s)
	}
	fmt.Fprintln(&output, au.Gray(8, ReplaceRuneAtIndex(border, '┴', 8)))

	return output.String()
}
