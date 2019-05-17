package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
		Aliases: []string{"new", "a"},
		Short:   "Add new snippet",
		Long:    `Add new snippet"`,
	}

	addCommandCmd = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command"},
		Short:   "Add new command snippet",
		Long:    `Add new command snippet`,
		Args:    cobra.NoArgs,
		Run:     addCommand,
	}

	addCodeCmd = &cobra.Command{
		Use:   "code CODE",
		Short: "Add new code snippet",
		Long:  `Add new code snippet`,
		Args:  cobra.NoArgs,
		Run:   addCode,
	}

	addBookmarkCmd = &cobra.Command{
		Use:     "url SNIPPET",
		Aliases: []string{"bookmark", "bm"},
		Short:   "Add new bookmark",
		Long:    `Add new bookmark`,
		Args:    cobra.MaximumNArgs(2),
		Run:     addBookmark,
	}
)

func addCommand(_ *cobra.Command, _ []string) {
	tags := tagsFlag
	content := ""
	if content = prompt("Command"); content == "" {
		fmt.Println("Canceled")
		return
	}

	description := ""
	if description = prompt("Description"); description == "" {
		fmt.Println("Canceled")
		return
	}

	command := snippet.NewCommand(content, description, tags)
	snippet.NewRepository().Add(command)

	fmt.Printf("Command added\n%s", command.String())
}

func addCode(_ *cobra.Command, args []string) {
	description := ""
	if description = prompt("Description"); description == "" {
		fmt.Println("Canceled")
		return
	}

	content := ""
	if content = fromEditor(); content == "" {
		fmt.Println("Canceled")
		return
	}

	language := prompt("Language")

	tags := tagsFlag

	code := snippet.NewCode(content, description, tags, language)
	snippet.NewRepository().Add(code)

	fmt.Printf("Code snippet added\n%s", code.String())
}

func addBookmark(_ *cobra.Command, args []string) {
	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	if url == "" {
		if url = prompt("URL"); url == "" {
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

func prompt(label string) string {
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

func fromEditor() string {
	fpath := os.TempDir() + "/gossip.tmp"
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	editor, _ := exec.LookPath(os.Getenv("EDITOR"))

	command := exec.Command(editor, fpath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err = command.Start(); err != nil {
		return ""
	}
	if err = command.Wait(); err != nil {
		return ""
	}

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return ""
	}
	return string(data)
}

func init() {
	addCmd.AddCommand(addCommandCmd)
	addCmd.AddCommand(addCodeCmd)
	addCmd.AddCommand(addBookmarkCmd)
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tagsFlag, "tags", "t", "", "tags")
}
