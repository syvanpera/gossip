package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var listCmd = &cobra.Command{
	Use:       "list [type]",
	Short:     "List snippets",
	Long:      `Lists snippets of given type, or all snippets if no type given`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"cmd", "code", "url"},
	Run:       list,
}

var tag, language string

func list(cmd *cobra.Command, args []string) {
	filters := snippet.Filters{Language: language}
	if tag != "" {
		filters.Tags = []string{tag}
	}
	if len(args) > 0 {
		filters.Type = args[0]
	}

	fmt.Println(filters)

	r := snippet.NewRepository()

	var snippets []snippet.Snippet
	snippets = r.FindWithFilters(filters)

	for _, s := range snippets {
		fmt.Println(s)
	}
}

func init() {
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", `Tag filter`)
	listCmd.Flags().StringVarP(&language, "language", "l", "", `Language filter`)
	rootCmd.AddCommand(listCmd)
}
