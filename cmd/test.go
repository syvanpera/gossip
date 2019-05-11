package cmd

import (
	"fmt"
	"io/ioutil"
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
	Run:   edit,
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Test template rendering",
	Run:   render,
}

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "Read from pipe",
	Run:   pipe,
}

var snippetTemplate = `
ID: {{.ID}}
Snippet: {{.Snippet}}
Description: {{.Description}}
Tags: {{range .Tags}}{{.}} {{end}}
Type: {{.Type}}{{if .Language}}Language: {{.Language}}{{end}}
`

func test(cmd *cobra.Command, args []string) {
}

func edit(cmd *cobra.Command, args []string) {
	fpath := os.TempDir() + "/gossip.tmp"
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	command := exec.Command("nvim", fpath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Start()
	if err != nil {
		log.Printf("Error while launching editor. Error: %v\n", err)
		log.Fatal(err)
	}
	err = command.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Printf("Error while reading. Error: %v\n", err)
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func render(cmd *cobra.Command, args []string) {
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
}

func pipe(cmd *cobra.Command, args []string) {
	// info, err := os.Stdin.Stat()
	// if err != nil {
	// 	panic(err)
	// }

	// if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
	// 	fmt.Println("The command is intended to work with pipes.")
	// 	fmt.Println("Usage: fortune | gocowsay")
	// 	return
	// }

	data, _ := ioutil.ReadAll(os.Stdin)
	fmt.Printf("stdin data: %v\n", string(data))

	// reader := bufio.NewReader(os.Stdin)
	// var output []rune

	// for {
	// 	input, _, err := reader.ReadRune()
	// 	if err != nil && err == io.EOF {
	// 		break
	// 	}
	// 	output = append(output, input)
	// }

	// for j := 0; j < len(output); j++ {
	// 	fmt.Printf("%c", output[j])
	// }
}
func init() {
	testCmd.AddCommand(editCmd)
	testCmd.AddCommand(renderCmd)
	testCmd.AddCommand(pipeCmd)
	rootCmd.AddCommand(testCmd)
}
