package bookmarks

import (
	"fmt"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

type Bookmark struct {
	ID uint `gorm:"primaryKey" json:"id"`

	URL         string `gorm:"unique" json:"url"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Flags       int    `json:"flags"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Bookmarks is used to present a list of bookmarks in JSON
type Bookmarks struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

// BookmarkCreateUpdate is the entity which is accepted for creates and updates to bookmarks
type BookmarkCreateUpdate struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Flags       int    `json:"flags"`
}

func (b Bookmark) DisplayText(plural bool) string {
	if plural {
		return "bookmarks"
	} else {
		return "bookmark"
	}
}

func (b Bookmark) Render(compact bool) string {
	var output strings.Builder

	if compact {
		fmt.Fprintf(&output, "%d|%s|%s|%s", b.ID, b.URL, b.Description, b.Tags)
	} else {
		colors := viper.GetBool("config.color") != viper.GetBool("color")
		au := aurora.NewAurora(colors)

		fmt.Fprintf(&output, "%s ", au.BrightCyan(fmt.Sprintf("%d.", b.ID)))
		fmt.Fprintln(&output, au.Bold(au.BrightGreen(b.Description)))
		fmt.Fprintf(&output, "   %s %s", au.BrightRed(">"), au.BrightYellow(b.URL))
		if b.Tags != "" {
			fmt.Fprintf(&output, "\n   %s %s", au.BrightRed("#"), au.BrightBlue(b.Tags))
			// fmt.Fprintf(&output, "\n   ")
			// for _, t := range strings.Split(b.Tags, ",") {
			// 	fmt.Fprintf(&output, "%s%s ", au.BrightRed("#"), au.BrightBlue(t))
			// }
		}
		// fmt.Fprintf(&output, "\n   %v %v", au.BrightMagenta(b.CreatedAt.Local()), au.BrightMagenta(b.UpdatedAt.Local()))
	}

	return output.String()
}
