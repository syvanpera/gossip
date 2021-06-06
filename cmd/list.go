package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/gossip"
	"github.com/syvanpera/gossip/ui"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List snippets",
		Long:    `List snippets`,
		Run:     listPrompt,
	}

	listCommandCmd = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command", "c"},
		Short:   "List command snippets",
		Long:    `List command snippets`,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			list(gossip.COMMAND)
		},
	}

	listCodeCmd = &cobra.Command{
		Use:     "code",
		Aliases: []string{"d"},
		Short:   "List code snippets",
		Long:    `List code snippets`,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			list(gossip.CODE)
		},
	}

	listBookmarkCmd = &cobra.Command{
		Use:     "url",
		Aliases: []string{"u", "bookmark", "bm", "b"},
		Short:   "List bookmarks",
		Long:    `List bookmarks`,
		Run: func(_ *cobra.Command, _ []string) {
			list(gossip.BOOKMARK)
		},
	}
)

var filterTags string
var languageFilter string

func listPrompt(_ *cobra.Command, _ []string) {
	choice, err := ui.Choose("List what", []string{"Bookmarks", "Commands", "Code snippets"})
	if err != nil {
		return
	}

	switch choice {
	case "Bookmarks":
		list(gossip.BOOKMARK)
	case "Commands":
		list(gossip.COMMAND)
	case "Code snippets":
		list(gossip.CODE)
	}
}

func list(gossipType int) {
	// filters := snippet.Filters{
	// 	Language: languageFilter,
	// 	Type:     t,
	// }
	// if filterTags != "" {
	// 	filters.Tags = filterTags
	// }

	result, err := gossipService.List(gossipType, filterTags)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range result {
		fmt.Println(r.Render())
	}
}

func init() {
	listCmd.AddCommand(listCommandCmd)
	listCmd.AddCommand(listCodeCmd)
	listCmd.AddCommand(listBookmarkCmd)
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "Filter by tags (comma separated)")
	listCodeCmd.Flags().StringVarP(&language, "language", "l", "", "Filter by language")
}
