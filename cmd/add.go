package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/ui"
)

var tags string

var (
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a", "new", "n"},
		Short:   "Create a new snippet",
		Long:    `Create a new snippet"`,
		Args:    cobra.MaximumNArgs(2),
		Run:     add,
	}

	addCommandCmd = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command", "c"},
		Short:   "Create a new command snippet",
		Long:    `Create a new command snippet`,
		Args:    cobra.MaximumNArgs(2),
		Run:     addCommand,
	}

	addCodeCmd = &cobra.Command{
		Use:     "code",
		Aliases: []string{"d"},
		Short:   "Create a new code snippet",
		Long:    `Create a new code snippet`,
		Args:    cobra.MaximumNArgs(1),
		Run:     addCode,
	}

	addBookmarkCmd = &cobra.Command{
		Use:     "url",
		Aliases: []string{"u", "bookmark", "bm", "b"},
		Short:   "Create a new bookmark",
		Long:    `Create a new bookmark`,
		Args:    cobra.MaximumNArgs(2),
		Run:     addBookmark,
	}
)

var service = snippet.NewService(snippet.NewSQLiteRepository())

func add(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		if matched, _ := regexp.MatchString("^https?://*", args[0]); matched {
			addBookmark(cmd, args)
			return
		}
	}

	choice, err := ui.Choose("Create what", []string{"Bookmark", "Command", "Code snippet"})
	if err != nil {
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

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	s, err := service.AddBookmark(content, description, tags)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Bookmark added\n%s", s.String())
}

func addCommand(_ *cobra.Command, args []string) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	s, err := service.AddCommand(content, description, tags)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Command added\n%s", s.String())
}

func addCode(_ *cobra.Command, args []string) {
	description := ""
	if len(args) > 0 {
		description = args[0]
	}

	s, err := service.AddCode("", description, tags)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Code snippet added\n%s", s.String())
}

func resolveContent(label string, args []string) string {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	if content == "" {
		var err error
		content, err = snippet.Edit(label, content)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	return content
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "tags")
}
