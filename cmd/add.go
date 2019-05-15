package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/util"
)

var tags string

var (
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"new"},
		Short:   "Add new snippet",
		Long:    `Add new snippet"`,
	}

	addCommandCmd = &cobra.Command{
		Use:     "command SNIPPET",
		Aliases: []string{"cmd"},
		Short:   "Add new command snippet",
		Long:    `Add new command snippet`,
		Args:    cobra.MinimumNArgs(1),
		Run:     addCommand,
	}

	addCodeCmd = &cobra.Command{
		Use:   "code SNIPPET",
		Short: "Add new code snippet",
		Long:  `Add new code snippet`,
		Args:  cobra.MinimumNArgs(1),
		Run:   addCode,
	}

	addBookmarkCmd = &cobra.Command{
		Use:     "bookmark SNIPPET",
		Aliases: []string{"bm"},
		Short:   "Add new bookmark",
		Long:    `Add new bookmark`,
		Args:    cobra.MinimumNArgs(1),
		Run:     addBookmark,
	}
)

func addCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("addCommand: %v", args)
}

func addCode(cmd *cobra.Command, args []string) {
	fmt.Printf("addCode: %v", args)
}

func addBookmark(cmd *cobra.Command, args []string) {
	fmt.Printf("addBookmark: %v %v\n", args, tags)
	fmt.Printf("tags: %v\n", strings.Split(tags, ","))

	url := args[0]
	description := util.ExtractTitleFromURL(url)
	fmt.Println("Got title [" + description + "]")

	bookmark := snippet.NewBookmark(url, description, tags)
	snippets := snippet.NewRepository()
	snippets.Upsert(bookmark)
	// bookmark := snippet.SnippetData{
	// 	Content:     url,
	// 	Description: description,
	// 	Type:        "BOOKMARK",
	// 	Tags:        strings.Split(tags, ","),
	// }
	// snippets.New(&bookmark)
	fmt.Printf("Added new bookmark\n%s", bookmark.String())
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "tags")
}
