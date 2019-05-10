package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
	"golang.org/x/crypto/ssh/terminal"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a snippet",
	Long:  `Show a snippet`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: show,
}

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func show(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("Snippet with ID %d not found\n", id)
		return
	}

	width, _, err := terminal.GetSize(0)
	if err != nil {
		width = 80
	}
	var sb strings.Builder
	borderStart := strings.Repeat("─", 8)
	borderEnd := strings.Repeat("─", width-(8+1))
	fmt.Fprintf(&sb, "%s┬%s\n", borderStart, borderEnd)
	fmt.Fprintf(&sb, "%s│ %s\n", center(fmt.Sprintf("#%d", s.ID), 8), s.Description)
	fmt.Fprintf(&sb, "%s┼%s\n", borderStart, borderEnd)

	for i, s := range strings.Split(s.Snippet, "\n") {
		fmt.Fprintf(&sb, center(strconv.Itoa(i+1), 8))
		fmt.Fprintf(&sb, "│ %s\n", s)
	}
	fmt.Fprintf(&sb, "%s┴%s\n", borderStart, borderEnd)

	// lexer := lexers.Analyse(s.Snippet)
	// fmt.Println(lexer.Config().Name)
	// quick.Highlight(os.Stdout, s.Snippet, lexer.Config().Name, "terminal16m", "dracula")

	fmt.Printf(sb.String())
}

func init() {
	rootCmd.AddCommand(showCmd)
}
