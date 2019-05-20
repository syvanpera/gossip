package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/ui"
)

var tagsFlag string

var (
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a", "new", "n"},
		Short:   "Add a new snippet",
		Long:    `Add a new snippet"`,
		Args:    cobra.MaximumNArgs(2),
		Run:     add,
	}

	addCommandCmd = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command", "c"},
		Short:   "Add a new command snippet",
		Long:    `Add a new command snippet`,
		Args:    cobra.MaximumNArgs(2),
		Run:     addCommand,
	}

	addCodeCmd = &cobra.Command{
		Use:     "code",
		Aliases: []string{"d"},
		Short:   "Add a new code snippet",
		Long:    `Add a new code snippet`,
		Args:    cobra.MaximumNArgs(1),
		Run:     addCode,
	}

	addBookmarkCmd = &cobra.Command{
		Use:     "url",
		Aliases: []string{"u", "bookmark", "bm", "b"},
		Short:   "Add a new bookmark",
		Long:    `Add a new bookmark`,
		Args:    cobra.MaximumNArgs(2),
		Run:     addBookmark,
	}
)

func add(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		if matched, _ := regexp.MatchString("^https?://*", args[0]); matched {
			addBookmark(cmd, args)
			return
		}
	}

	choice, err := ui.Choose("Add what", []string{"Bookmark", "Command", "Code snippet"})

	if err != nil {
		fmt.Println("Canceled")
		return
	}

	switch choice {
	case "Bookmark":
		addBookmark(cmd, args)
	case "Command":
		addCommand(cmd, args)
	case "Code snippet":
		addCode(cmd, args)
	}
}

func addBookmark(_ *cobra.Command, args []string) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	if content == "" {
		if content = ui.Prompt("URL", ""); content == "" {
			fmt.Println("Canceled")
			return
		}
	}

	if matched, _ := regexp.MatchString("^https?://*", content); !matched {
		content = "http://" + content
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	tags := tagsFlag

	if description == "" || tags == "" {
		if meta := meta.Extract(content); meta != nil {
			if description == "" {
				description = meta.Description
			}
			if tags == "" {
				tags = meta.Tags
			}
		}
	}

	bookmark := snippet.NewBookmark(content, description, tags)
	snippet.NewRepository().Add(&bookmark)

	fmt.Printf("Bookmark added\n%s", bookmark.String())
}

func addCommand(_ *cobra.Command, args []string) {
	tags := tagsFlag
	content := ""
	if len(args) > 0 {
		content = args[0]
	}
	// TODO Refactor to use the Edit method from the Snippet
	if content == "" {
		if content = ui.Prompt("Command", ""); content == "" {
			fmt.Println("Canceled")
			return
		}
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}
	if description == "" {
		if description = ui.Prompt("Description", ""); description == "" {
			fmt.Println("Canceled")
			return
		}
	}

	command := snippet.NewCommand(content, description, tags)
	snippet.NewRepository().Add(&command)

	fmt.Printf("Command added\n%s", command.String())
}

func addCode(_ *cobra.Command, args []string) {
	description := ""
	if len(args) > 0 {
		description = args[0]
	}
	if description == "" {
		if description = ui.Prompt("Description", ""); description == "" {
			fmt.Println("Canceled")
			return
		}
	}

	content := ""
	if content = ui.Editor(""); content == "" {
		fmt.Println("Canceled")
		return
	}

	language := ui.Prompt("Language", "")

	tags := tagsFlag

	code := snippet.NewCode(content, description, tags, language)
	snippet.NewRepository().Add(&code)

	fmt.Printf("Code snippet added\n%s", code.String())
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tagsFlag, "tags", "t", "", "tags")
}
