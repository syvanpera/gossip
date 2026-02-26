// cmd/search.go
package cmd

import (
	"fmt"
	"os"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search [query]",
	Short:   "Search bookmarks by keyword or tag",
	Aliases: []string{"s", "find", "f"},
	Args:    cobra.ExactArgs(1), // Requires exactly one argument: the search term
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		storagePath := config.GetStoragePath()

		results, err := storage.Search(storagePath, query)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✗ Failed to search:"), err)
			os.Exit(1)
		}

		// Handle the case where nothing is found
		if len(results) == 0 {
			fmt.Printf("No bookmarks found matching '%s'.\n", ui.StyleTitle.Render(query))
			return
		}

		// Print a nice header, then the results
		fmt.Printf("%s Found %d result(s) for '%s':\n\n", ui.StyleSuccess.Render("✓"), len(results), ui.StyleTitle.Render(query))
		ui.PrintSearchResults(results, query)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
