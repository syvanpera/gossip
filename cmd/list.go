package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var listCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List snippets",
	Long:  `Lists snippets of given type, or all snippets if no type given`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return cobra.OnlyValidArgs(cmd, args)
	},
	Run: list,
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

	r := snippet.NewRepository()
	snippets := r.FindWithFilters(filters)

	for _, s := range snippets {
		fmt.Println(s)
	}
}

func init() {
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", `Tag filter`)
	listCmd.Flags().StringVarP(&language, "language", "l", "", `Language filter`)
	rootCmd.AddCommand(listCmd)
}
