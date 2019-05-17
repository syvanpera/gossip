package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/ui"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List snippets",
		Long:    `List snippets`,
		Run:     listDefault,
	}

	listCommandCmd = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command"},
		Short:   "List command snippets",
		Long:    `List command snippets`,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			list(snippet.COMMAND)
		},
	}

	listCodeCmd = &cobra.Command{
		Use:   "code",
		Short: "List code snippets",
		Long:  `List code snippets`,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			list(snippet.CODE)
		},
	}

	listBookmarkCmd = &cobra.Command{
		Use:     "url",
		Aliases: []string{"u", "bookmark", "bm", "b"},
		Short:   "List bookmarks",
		Long:    `List bookmarks`,
		Run: func(_ *cobra.Command, _ []string) {
			list(snippet.BOOKMARK)
		},
	}
)

var tags, language string

func listDefault(_ *cobra.Command, _ []string) {
	choice, err := ui.Choose("List what", []string{"Bookmarks", "Commands", "Code snippets"})

	if err != nil {
		fmt.Println("Canceled")
		return
	}

	switch choice {
	case "Bookmarks":
		list(snippet.BOOKMARK)
	case "Commands":
		list(snippet.COMMAND)
	case "Code snippets":
		list(snippet.CODE)
	}
}

func list(t snippet.SnippetType) {
	filters := snippet.Filters{
		Language: language,
		Type:     t,
	}
	if tags != "" {
		filters.Tags = tags
	}

	r := snippet.NewRepository()
	snippets := r.FindWithFilters(filters)

	for _, s := range snippets {
		fmt.Println(s)
	}
}

func init() {
	listCmd.AddCommand(listCommandCmd)
	listCmd.AddCommand(listCodeCmd)
	listCmd.AddCommand(listBookmarkCmd)
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "Tags filter (comma separated)")
	listCodeCmd.Flags().StringVarP(&language, "language", "l", "", "Language filter")
}
