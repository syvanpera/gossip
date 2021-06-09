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

func list() {
	// filters := snippet.Filters{
	// 	Language: languageFilter,
	// 	Type:     t,
	// }
	// if filterTags != "" {
	// 	filters.Tags = filterTags
	// }

	result, err := gossipService.List(tagsFilter)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range result {
		fmt.Println(r.Render())
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&tagsFilter, "tags", "t", "", "Filter by tags (comma separated)")
}
