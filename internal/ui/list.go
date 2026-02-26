// internal/ui/list.go
package ui

import (
	"fmt"
	"strings"

	"github.com/syvanpera/gossip/internal/storage"
)

// PrintBookmarks iterates through the bookmarks and prints them with Lip Gloss styles
func PrintBookmarks(bookmarks []storage.Bookmark) {
	if len(bookmarks) == 0 {
		fmt.Println(StyleTitle.Render("No bookmarks found. Try adding one!"))
		return
	}

	for _, b := range bookmarks {
		// Line 1: ID. Title
		// Appends a period to the ID to match the numbered list style
		idStr := StyleID.Render(fmt.Sprintf("%s.", b.ID))
		fmt.Printf("%s %s\n", idStr, StyleTitle.Render(b.Title))

		// Line 2:   > URL
		// Using the existing red StyleFailed for the prefix to match the image's accent color
		prefixURL := StyleFailed.Render("  >")
		fmt.Printf("%s %s\n", prefixURL, StyleURL.Render(b.URL))

		// Line 3:   # tag1,tag2,tag3 (if any exist)
		if len(b.Tags) > 0 {
			prefixTags := StyleFailed.Render("  #")
			// Join tags with a comma and no spaces to match the screenshot
			tagsStr := StyleTags.Render(strings.Join(b.Tags, ","))
			fmt.Printf("%s %s\n", prefixTags, tagsStr)
		}

		// Line 4:   // comment (if it exists, so we don't lose description data)
		if b.Comment != "" {
			prefixComment := StyleID.Render("  //")
			fmt.Printf("%s %s\n", prefixComment, StyleComment.Render(b.Comment))
		}

		fmt.Println() // Add a blank line between entries
	}
}
