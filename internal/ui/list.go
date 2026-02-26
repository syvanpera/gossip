// internal/ui/list.go
package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/syvanpera/gossip/internal/storage"
)

// PrintBookmarks iterates through the bookmarks and prints them with Lip Gloss styles
func PrintBookmarks(bookmarks []storage.Bookmark) {
	if len(bookmarks) == 0 {
		fmt.Println(StyleTitle.Render("No bookmarks found. Try adding one!"))
		return
	}

	fmt.Println() // Add a blank line before the list for better spacing

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

// PrintSearchResults highlights the search query in the output
func PrintSearchResults(bookmarks []storage.Bookmark, query string) {
	if len(bookmarks) == 0 {
		return
	}

	for _, b := range bookmarks {
		// Highlight title
		idStr := StyleID.Render(fmt.Sprintf("%s.", b.ID))
		titleStr := highlightMatch(b.Title, query, StyleTitle)
		fmt.Printf("%s %s\n", idStr, titleStr)

		// Highlight URL
		prefixURL := StyleFailed.Render("  >")
		urlStr := highlightMatch(b.URL, query, StyleURL)
		fmt.Printf("%s %s\n", prefixURL, urlStr)

		// Highlight Tags
		if len(b.Tags) > 0 {
			prefixTags := StyleFailed.Render("  #")

			var highlightedTags []string
			for _, tag := range b.Tags {
				highlightedTags = append(highlightedTags, highlightMatch(tag, query, StyleTags))
			}

			// Join tags, ensuring the commas get the default tag style
			tagsStr := strings.Join(highlightedTags, StyleTags.Render(","))
			fmt.Printf("%s %s\n", prefixTags, tagsStr)
		}

		// Highlight Comment
		if b.Comment != "" {
			prefixComment := StyleID.Render("  //")
			commentStr := highlightMatch(b.Comment, query, StyleComment)
			fmt.Printf("%s %s\n", prefixComment, commentStr)
		}

		fmt.Println()
	}
}

// highlightMatch applies the base style, but highlights the matching query safely
func highlightMatch(text, query string, baseStyle lipgloss.Style) string {
	if query == "" {
		return baseStyle.Render(text)
	}

	// Create a case-insensitive regex for the exact query
	re := regexp.MustCompile(`(?i)(` + regexp.QuoteMeta(query) + `)`)
	indices := re.FindAllStringIndex(text, -1)

	if len(indices) == 0 {
		return baseStyle.Render(text) // No match, just render normally
	}

	var result strings.Builder
	lastIdx := 0

	for _, match := range indices {
		start, end := match[0], match[1]

		// Add unstyled text before the match (styled with baseStyle)
		if start > lastIdx {
			result.WriteString(baseStyle.Render(text[lastIdx:start]))
		}

		// Add matched text (styled with highlightStyle)
		result.WriteString(StyleHighlight.Render(text[start:end]))
		lastIdx = end
	}

	// Add any remaining text after the final match
	if lastIdx < len(text) {
		result.WriteString(baseStyle.Render(text[lastIdx:]))
	}

	return result.String()
}
