// cmd/clip.go
package cmd

import (
	"fmt"
	"os"
	"strings" // For autocomplete

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var clipCmd = &cobra.Command{
	Use:   "clip [id]",
	Short: "Copy a bookmark's URL to your clipboard",
	Args:  cobra.ExactArgs(1), // Requires exactly one argument: the ID
	// Dynamic TAB completion for the ID
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		storagePath := config.GetStoragePath()
		bookmarks, err := storage.Load(storagePath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, b := range bookmarks {
			if strings.HasPrefix(b.ID, toComplete) {
				completion := fmt.Sprintf("%s\t%s", b.ID, b.Title)
				completions = append(completions, completion)
			}
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		storagePath := config.GetStoragePath()

		// 1. Fetch the bookmark
		bm, err := storage.GetByID(storagePath, id)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✖ Error:"), err)
			os.Exit(1)
		}

		// 2. Copy the URL to the system clipboard
		err = clipboard.WriteAll(bm.URL)
		if err != nil {
			fmt.Printf("\n%s Failed to copy to clipboard: %v\n", ui.StyleFailed.Render("✗"), err)
			os.Exit(1)
		}

		// 3. Print a success message
		fmt.Printf("\n%s URL for '%s' copied to clipboard!\n", ui.StyleSuccess.Render("✓"), ui.StyleTitle.Render(bm.Title))
	},
}

func init() {
	rootCmd.AddCommand(clipCmd)
}
