package cmd

import (
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing",
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Test Edit",
	RunE:  edit,
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Test template rendering",
	RunE:  render,
}

var snippetTemplate = `
ID: {{.ID}}
Snippet: {{.Snippet}}
Description: {{.Description}}
Tags: {{range .Tags}}{{.}} {{end}}
Type: {{.Type}}{{if .Language}}Language: {{.Language}}{{end}}
`

func test(cmd *cobra.Command, args []string) error {
	return nil
}

func edit(cmd *cobra.Command, args []string) error {
	fpath := os.TempDir() + "/gossip.tmp"
	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("1")
		log.Fatal(err)
	}
	f.Close()

	command := exec.Command("nvim", fpath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Start()
	if err != nil {
		log.Printf("2")
		log.Fatal(err)
	}
	err = command.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	} else {
		log.Printf("Successfully edited.")
	}

	return nil
}

func render(cmd *cobra.Command, args []string) error {
	snippet := snippet.Snippet{
		ID:          1,
		Snippet:     "Test snippet",
		Description: "Test description",
		Tags:        []string{"tag1", "tag2"},
		Type:        "code",
		Language:    "go",
	}

	t, err := template.New("snippet").Parse(snippetTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, snippet)
	if err != nil {
		panic(err)
	}

	return nil
}

func init() {
	testCmd.AddCommand(editCmd)
	testCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(testCmd)
}
