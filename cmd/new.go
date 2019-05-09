package cmd

import (
	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add new snippet",
	Long:  `Add new snippet`,
	RunE:  new,
}

func new(cmd *cobra.Command, args []string) error {
	snippets := snippet.NewRepository()
	snippet := snippet.Snippet{
		Snippet:     "ls -la",
		Description: snippet.ToNullString("Lists the contents of a directory"),
		Type:        "cmd",
	}
	snippets.New(&snippet)

	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
