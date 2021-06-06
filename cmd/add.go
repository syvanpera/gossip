package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/gossip"
	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/ui"
	"github.com/syvanpera/gossip/util"
)

var tags string
var stdin bool
var language string

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
		Use:     "command",
		Aliases: []string{"cmd", "c"},
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
		Use:     "bookmark",
		Aliases: []string{"b"},
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
	content := ""
	if len(args) > 0 {
		content = args[0]
	} else if content = ui.Prompt("URL", content); content == "" {
		return
	}

	if matched, _ := regexp.MatchString("^https?://*", content); !matched {
		content = "https://" + content
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	if description == "" || tags == "" {
		if meta := meta.Extract(content); meta != nil {
			if description == "" {
				description = meta.Description
			}
			if tags == "" {
				tags = fmt.Sprintf("%s,%s", tags, meta.Tags)
			}
		}
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return
		}
	}

	s, err := gossipService.Create(gossip.BOOKMARK, content, description, tags)
	if err != nil {
		util.PrintError(err)
		return
	}

	fmt.Println(s.Render())
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

	s, err := service.CreateCommand(content, description, tags)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Command added\n%s", s.Render())
}

func addCode(_ *cobra.Command, args []string) {
	quiet := viper.GetBool("quiet")
	description := ""
	if len(args) > 0 {
		description = args[0]
	}

	s, err := service.CreateCode("", description, tags, language)
	if err != nil {
		if !quiet {
			fmt.Println(err)
		}
		return
	}

	if !quiet {
		fmt.Printf("Code snippet added\n%s", s.Render())
	}
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "tags")
	addCodeCmd.Flags().BoolVarP(&stdin, "stdin", "s", false, "read from stdin")
	addCodeCmd.Flags().StringVarP(&language, "language", "l", "", "language")
}
