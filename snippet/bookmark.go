package snippet

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/ui"
)

type Bookmark struct {
	data SnippetData
}

func (b *Bookmark) Type() SnippetType  { return BOOKMARK }
func (b *Bookmark) Data() *SnippetData { return &b.data }

func (b *Bookmark) Execute() error {
	br := viper.GetString("config.browser")
	fmt.Printf("Okay, opening link in %s browser...\n", br)
	url := b.data.Content
	if matched, _ := regexp.MatchString("^http(s)?://*", url); !matched {
		url = "http://" + url
	}
	if br == "default" {
		browser.OpenURL(url)
	} else {
		openInBrowser(br, url)
	}

	return nil
}

func (b *Bookmark) Edit() error {
	content := b.Data().Content
	if content = ui.Prompt("URL", content); content == "" {
		return ErrEditCanceled
	}
	b.Data().Content = content

	description := b.Data().Description
	if description = ui.Prompt("Description", b.Data().Description); description == "" {
		return ErrEditCanceled
	}
	b.Data().Description = description

	return nil
}

func (b *Bookmark) String() string {
	colors := viper.GetBool("config.color") != viper.GetBool("color")
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

func NewBookmark(url, description string, tags string) Bookmark {
	bookmark := Bookmark{
		data: SnippetData{
			Content:     url,
			Description: description,
			Tags:        strings.ToLower(tags),
			Type:        BOOKMARK,
		},
	}

	return bookmark
}

func openInBrowser(b, url string) {
	var command *exec.Cmd
	command = exec.Command(b, url)
	command.Start()
}
