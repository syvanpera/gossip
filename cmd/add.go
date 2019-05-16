package cmd

import (
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
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
		Use:     "cmd COMMAND",
		Aliases: []string{"command"},
		Short:   "Add new command snippet",
		Long:    `Add new command snippet`,
		Args:    cobra.MinimumNArgs(1),
		Run:     addCommand,
	}

	addCodeCmd = &cobra.Command{
		Use:   "code CODE",
		Short: "Add new code snippet",
		Long:  `Add new code snippet`,
		Args:  cobra.MinimumNArgs(1),
		Run:   addCode,
	}

	addBookmarkCmd = &cobra.Command{
		Use:     "bm SNIPPET",
		Aliases: []string{"bookmark", "url"},
		Short:   "Add new bookmark",
		Long:    `Add new bookmark`,
		// Args:    cobra.MinimumNArgs(1),
		Run: addBookmark,
	}
)

func addCommand(cmd *cobra.Command, args []string) {
	command := args[0]
	fmt.Printf("addCommand: %v\n", command)
}

func addCode(cmd *cobra.Command, args []string) {
	fmt.Printf("addCode: %v", args)
}

func addBookmark(cmd *cobra.Command, args []string) {
	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	if url == "" {
		if url = promptFor("URL"); url == "" {
			fmt.Println("Canceled")
			return
		}
	}

	if matched, _ := regexp.MatchString("^https?://*", url); !matched {
		url = "http://" + url
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	tags := tagsFlag

	if description == "" || tags == "" {
		if meta := meta.Extract(url); meta != nil {
			if description == "" {
				description = meta.Description
			}
			if tags == "" {
				tags = meta.Tags
			}
		}
	}

	bookmark := snippet.NewBookmark(url, description, tags)
	snippet.NewRepository().Add(bookmark)

	fmt.Printf("Bookmark added\n%s", bookmark.String())
}

func promptFor(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	input, err := prompt.Run()
	if err != nil {
		return ""
	}
	if input == "" {
		return ""
	}

	return input
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tagsFlag, "tags", "t", "", "tags")
}
