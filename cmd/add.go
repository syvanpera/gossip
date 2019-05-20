package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/ui"
)

var tags string

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
	data := snippet.SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        snippet.COMMAND,
	}

	data.Content = resolveContent("URL", args)
	if data.Content == "" {
		return
	}

	if matched, _ := regexp.MatchString("^https?://*", data.Content); !matched {
		data.Content = "http://" + data.Content
	}

	if len(args) > 1 {
		data.Description = args[1]
	}
	if data.Description == "" || data.Tags == "" {
		if meta := meta.Extract(data.Content); meta != nil {
			if data.Description == "" {
				data.Description = meta.Description
			}
			if data.Tags == "" {
				data.Tags = meta.Tags
			}
		}
	}

	s := snippet.New(data)
	snippet.NewRepository().Add(s)

	fmt.Printf("Bookmark added\n%s", s.String())
}

func addCommand(_ *cobra.Command, args []string) {
	data := snippet.SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        snippet.COMMAND,
	}

	data.Content = resolveContent("Command", args)
	if data.Content == "" {
		return
	}

	if len(args) > 1 {
		data.Description = args[1]
	}
	if data.Description == "" {
		var err error
		data.Description, err = snippet.Edit("Description", data.Description)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	s := snippet.New(data)
	snippet.NewRepository().Add(s)

	fmt.Printf("Command added\n%s", s.String())
}

func addCode(_ *cobra.Command, args []string) {
	data := snippet.SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        snippet.CODE,
	}

	if len(args) > 0 {
		data.Description = args[0]
	}
	if data.Description == "" {
		if data.Description = ui.Prompt("Description", ""); data.Description == "" {
			fmt.Println("Canceled")
			return
		}
	}

	if data.Content = ui.Editor(""); data.Content == "" {
		fmt.Println("Canceled")
		return
	}

	data.Language = ui.Prompt("Language", "")

	s := snippet.New(data)
	snippet.NewRepository().Add(s)

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
