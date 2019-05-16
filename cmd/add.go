package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/snippet"
)

var tagsFlag string

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
	url := args[0]
	if matched, _ := regexp.MatchString("^https?://*", url); !matched {
		url = "http://" + url
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	var tags []string
	if tagsFlag != "" {
		tags = strings.Split(tagsFlag, ",")
	}

	if description == "" || len(tags) == 0 {
		if meta := meta.Extract(url); meta != nil {
			description = meta.Description
			tags = append(tags, meta.Tags...)
		}
	}

	bookmark := snippet.NewBookmark(url, description, tags)
	if err := snippet.NewRepository().Add(bookmark); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("New bookmark added\n%s", bookmark.String())
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tagsFlag, "tags", "t", "", "tags")
}
