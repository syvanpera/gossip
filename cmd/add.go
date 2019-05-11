package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/logrusorgru/aurora"
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
	colors := viper.GetBool("colorOutput") != viper.GetBool("color")
	au := aurora.NewAurora(colors)
	snippet, err := scan(au.Yellow("Snippet> ").String())
	if err == errCanceled {
		return
	}
	fmt.Println(snippet, err)
	description, err := scan(au.Cyan("Description> ").String())
	if err == errCanceled {
		return
	}
	fmt.Println(description, err)
	// snippets := snippet.NewRepository()
	// s := snippet.Snippet{
	// 	Snippet:     "ls -la",
	// 	Description: "Lists the contents of a directory",
	// 	Type:        "cmd",
	// }
	// snippets.New(&s)
	// fmt.Printf("Added new snippet\n%s", s)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
