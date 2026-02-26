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

	fmt.Println() // Add a blank line at the start

	for _, b := range bookmarks {
		// Format the tags like [tag1, tag2]
		tagStr := ""
		if len(b.Tags) > 0 {
			tagStr = StyleTags.Render(fmt.Sprintf("[%s]", strings.Join(b.Tags, ", ")))
		}

		// Print the ID and Title
		fmt.Printf("%s %s %s\n", StyleID.Render(b.ID), StyleTitle.Render(b.Title), tagStr)

		// Print the URL
		fmt.Printf("  %s\n", StyleURL.Render(b.URL))

		// Print the comment/description if it exists
		if b.Comment != "" {
			fmt.Printf("  %s\n", StyleComment.Render(b.Comment))
		}

		fmt.Println() // Add a blank line between entries
	}
}
