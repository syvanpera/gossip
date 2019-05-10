package cmd

import (
	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add new snippet",
	Long:  `Add new snippet`,
	Run:   new,
}

func new(cmd *cobra.Command, args []string) {
	snippets := snippet.NewRepository()
	snippet := snippet.Snippet{
		Snippet: `func show(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("Snippet with ID %d not found\n", id)
		return
	}
	// fmt.Println(s)

	width, _, err := terminal.GetSize(0)
	if err != nil {
		width = 80
	}
	var sb strings.Builder
	borderStart := strings.Repeat("─", 8)
	borderEnd := strings.Repeat("─", width-(8+1))
	fmt.Fprintf(&sb, "%s┬%s\n", borderStart, borderEnd)
	fmt.Fprintf(&sb, center(fmt.Sprintf("#%d", s.ID), 8))
	fmt.Fprintf(&sb, "│\n%s┼%s\n", borderStart, borderEnd)

	for i, s := range strings.Split(s.Snippet, "\n") {
		fmt.Fprintf(&sb, center(strconv.Itoa(i+1), 8))
		fmt.Fprintf(&sb, "│ %s\n", s)
		// fmt.Fprintf(&sb, "%8d│ %s\n", i+1, s)
	}
	fmt.Fprintf(&sb, "%s┴%s\n", borderStart, borderEnd)

	foo := "in the middle"
	w := 110 // or whatever

	fmt.Println(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(foo))/2, foo)))

	// lexer := lexers.Analyse(s.Snippet)
	// fmt.Println(lexer.Config().Name)
	// quick.Highlight(os.Stdout, s.Snippet, lexer.Config().Name, "terminal16m", "dracula")

	fmt.Println(sb.String())
}`,

		Description: "Some code snippet",
		Type:        "code",
		Language:    "go",
		Tags:        []string{"go", "example"},
	}
	snippets.New(&snippet)
}

func init() {
	rootCmd.AddCommand(newCmd)
}
