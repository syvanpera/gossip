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
	Run:   list,
}

var tag string

func list(cmd *cobra.Command, args []string) {
	// _type := ""
	// if len(args) > 0 {
	// _type = args[0]
	// }

	r := snippet.NewRepository()
	// fmt.Println(r.Get(1))

	var snippets []snippet.Snippet
	if tag != "" {
		snippets = r.FindWithTag(tag)
	} else {
		snippets = r.FindAll()
	}

	for _, s := range snippets {
		fmt.Println(s)
	}
}

func init() {
	listCmd.Flags().StringVarP(&tag, "tag", "t", "", `Tag filter`)
	rootCmd.AddCommand(listCmd)
}
