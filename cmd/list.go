// cmd/list.go
package cmd

import (
	"fmt"
	"os"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all saved bookmarks",
	Aliases: []string{"ls", "show"},
	Run: func(cmd *cobra.Command, args []string) {
		storagePath := config.GetStoragePath()

		// Load bookmarks from JSON
		bookmarks, err := storage.Load(storagePath)
		if err != nil {
			fmt.Println(ui.StyleFailed.Render("✗ Failed to load bookmarks:"), err)
			os.Exit(1)
		}

		// Print them beautifully
		ui.PrintBookmarks(bookmarks)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
