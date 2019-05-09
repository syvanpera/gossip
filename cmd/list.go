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
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	fmt.Println("list: ", args)
	// _type := ""
	// if len(args) > 0 {
	// _type = args[0]
	// }

	r := snippet.NewRepository()
	// fmt.Println(r.Get(1))
	snippets := r.FindAll()

	for _, s := range snippets {
		fmt.Println(s)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
