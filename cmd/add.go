package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"new"},
	Short:   "Add new snippet",
	Long:    `Add new snippet`,
	Run:     add,
}

func add(cmd *cobra.Command, args []string) {
	snippets := snippet.NewRepository()
	s := snippet.Snippet{
		Snippet:     "ls -la",
		Description: "Lists the contents of a directory",
		Type:        "cmd",
	}
	snippets.New(&s)
	fmt.Printf("Added new snippet\n%s", s)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
