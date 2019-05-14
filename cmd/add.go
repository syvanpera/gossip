package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"new"},
	Short:   "Add new snippet",
	Long:    `Add new snippet`,
	Run:     add,
}

var errCanceled = errors.New("Input canceled")

func scan(message string) (string, error) {
	tempFile := os.TempDir() + "/gossip.history"
	l, err := readline.NewEx(&readline.Config{
		Prompt:          message,
		HistoryFile:     tempFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", errCanceled
}

func add(cmd *cobra.Command, args []string) {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)
	// s, err := scan(au.Yellow("Snippet> ").String())
	// if err == errCanceled {
	// 	return
	// }
	// fmt.Println(s, err)
	description, err := scan(au.Cyan("Description> ").String())
	if err == errCanceled {
		return
	}
	fmt.Println(description)

	items := []string{"Go", "Javascript", "Elm", "Typescript"}
	index := -1
	var result string

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "Language",
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You chose %s\n", result)
	// snippets := snippet.NewRepository()
	// newSnippet := snippet.Snippet{
	// 	Snippet:     "ls -la",
	// 	Description: "Lists the contents of a directory",
	// 	Type:        "cmd",
	// }
	// snippets.New(&newSnippet)
	// fmt.Printf("Added new snippet\n%s", newSnippet)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
