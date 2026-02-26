// cmd/add.go
package cmd

import (
	"fmt"
	"os"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tags []string

var addCmd = &cobra.Command{
	Use:     "add [url]",
	Short:   "Add a new bookmark",
	Aliases: []string{"a", "new"},
	Args:    cobra.ExactArgs(1), // Requires exactly one argument: the URL
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		storagePath := config.GetStoragePath()

		// Initialize our Bubble Tea spinner model
		m := ui.NewAddModel(url, tags, storagePath)
		p := tea.NewProgram(m)

		// Run the UI program
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// Add a flag to allow users to tag their bookmarks
	addCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Comma-separated list of tags (e.g., -t go,cli,tutorial)")
}
