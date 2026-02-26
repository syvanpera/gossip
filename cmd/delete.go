// cmd/delete.go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a bookmark by its ID",
	Args:  cobra.ExactArgs(1),
	// This function powers the dynamic TAB completion
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// We only want to autocomplete the first argument (the ID)
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
			// If what the user typed so far matches the start of an ID
			if strings.HasPrefix(b.ID, toComplete) {
				// The \t separates the value from the description in Cobra
				completion := fmt.Sprintf("%s\t%s", b.ID, b.Title)
				completions = append(completions, completion)
			}
		}

		// Return the matches, and tell Cobra NOT to suggest local files
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		storagePath := config.GetStoragePath()

		err := storage.Delete(storagePath, id)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✖ Failed to delete:"), err)
			os.Exit(1)
		}

		fmt.Printf("%s Bookmark '%s' deleted successfully.\n", ui.StyleSuccess.Render("✔ Success!"), ui.StyleID.Render(id))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
