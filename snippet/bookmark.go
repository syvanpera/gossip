package snippet

import (
	"fmt"
	"os/exec"
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
	br := viper.GetString("config.browser")
	url := b.data.Content
	if matched, _ := regexp.MatchString("^http(s)?://*", url); !matched {
		url = "http://" + url
	}
	if _, err := exec.LookPath(br); err == nil {
		fmt.Printf("Okay, opening link in %s...\n", br)
		command := exec.Command(br, url)
		command.Start()
	} else {
		fmt.Println("Okay, opening link in default browser...")
		browser.OpenURL(url)
	}

	return nil
}

func (b *Bookmark) Edit() error {
	s, err := Edit("URL", b.Data().Content)
	if err != nil {
		return err
	}
	b.Data().Content = s

	s, err = Edit("Description", b.Data().Description)
	if err != nil {
		return err
	}
	b.Data().Description = s

	return nil
}

func (b *Bookmark) Render() string {
	colors := viper.GetBool("config.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder

	fmt.Fprintf(&output, "%s ", au.BrightCyan(fmt.Sprintf("%d.", b.data.ID)))
	fmt.Fprintln(&output, au.Bold(au.BrightGreen(b.data.Description)))
	fmt.Fprintf(&output, "   %s %s", au.BrightRed(">"), au.BrightYellow(b.data.Content))
	if b.data.Tags != "" {
		fmt.Fprintf(&output, "\n   %s %s", au.BrightRed("#"), au.BrightBlue(b.data.Tags))
	}

	return output.String()
}

func (b *Bookmark) String(plural bool) string {
	if plural {
		return "bookmarks"
	} else {
		return "bookmark"
	}
}
