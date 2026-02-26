// cmd/delete.go
package cmd

import (
	"fmt"
	"os"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete a bookmark by its ID",
	Aliases: []string{"del", "rm", "remove"},
	Args:    cobra.ExactArgs(1), // Requires exactly one argument: the ID
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		storagePath := config.GetStoragePath()

		err := storage.Delete(storagePath, id)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✗ Failed to delete:"), err)
			os.Exit(1)
		}

		fmt.Printf("%s Bookmark '%s' deleted successfully.\n", ui.StyleSuccess.Render("✓ Success!"), ui.StyleID.Render(id))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
