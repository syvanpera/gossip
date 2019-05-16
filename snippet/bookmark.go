package snippet

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

type Bookmark struct {
	data SnippetData
}

func (b *Bookmark) Type() SnippetType  { return BOOKMARK }
func (b *Bookmark) Data() *SnippetData { return &b.data }

func (b *Bookmark) Execute() error {
	fmt.Println("Okay, opening link in default browser...")
	url := b.data.Content
	if matched, _ := regexp.MatchString("^http(s)?://*", url); !matched {
		url = "http://" + url
	}
	browser.OpenURL(url)

	return nil
}

func (b *Bookmark) String() string {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder

	fmt.Fprintf(&output, "\n%s ", au.BrightCyan(fmt.Sprintf("%d.", b.data.ID)))
	fmt.Fprintln(&output, au.Bold(au.BrightGreen(b.data.Description)))
	fmt.Fprintf(&output, "   %s %s\n", au.BrightRed(">"), au.BrightYellow(b.data.Content))
	if b.data.Tags != "" {
		fmt.Fprintf(&output, "   %s %s\n", au.BrightRed("#"), au.BrightBlue(b.data.Tags))
	}

	return output.String()
}

func NewBookmark(url, description string, tags string) *Bookmark {
	bookmark := Bookmark{
		data: SnippetData{
			Content:     url,
			Description: description,
			Tags:        tags,
			Type:        BOOKMARK,
		},
	}

	return &bookmark
}
