package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List bookmarks",
		Long:    `List bookmarks`,
		Run: func(_ *cobra.Command, _ []string) {
			list()
		},
	}
)

var tagsFilter string
var compact bool

func list() {
	result, err := service.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(result) == 0 {
		fmt.Println("No bookmarks")
	}

	for _, r := range result {
		fmt.Printf("%s\n", r.Render(compact))
		if !compact {
			fmt.Println()
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&tagsFilter, "tags", "t", "", "Filter by tags (comma separated)")
	listCmd.PersistentFlags().BoolVarP(&compact, "compact", "c", false, "Compact one line output (useful for dmenu etc)")
}
