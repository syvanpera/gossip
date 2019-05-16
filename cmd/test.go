package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/chzyer/readline"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/snippet"
)

var testErrCanceled = errors.New("Input canceled")

var (
	testCmd = &cobra.Command{
		Use:   "test",
		Short: "Testing",
	}

	testEditCmd = &cobra.Command{
		Use:   "edit",
		Short: "Test Edit",
		Run:   testEdit,
	}

	testRenderCmd = &cobra.Command{
		Use:   "render",
		Short: "Test template rendering",
		Run:   testRender,
	}

	testPipeCmd = &cobra.Command{
		Use:   "pipe",
		Short: "Read from pipe",
		Run:   testPipe,
	}

	testLanguagesCmd = &cobra.Command{
		Use:   "languages",
		Short: "List supported languages",
		Run:   testLanguages,
	}

	testAddCmd = &cobra.Command{
		Use:     "add [type]",
		Aliases: []string{"new"},
		Short:   "Add new snippet",
		Long:    `Add new snippet of type "type"`,
		Run:     testAdd,
	}
)

var snippetTemplate = `
ID: {{.ID}}
Snippet: {{.Snippet}}
Description: {{.Description}}
Tags: {{range .Tags}}{{.}} {{end}}
Type: {{.Type}}{{if .Language}}Language: {{.Language}}{{end}}
`

func test(cmd *cobra.Command, args []string) {
}

func testEdit(cmd *cobra.Command, args []string) {
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
		log.Printf("Error while launching editor. Error: %v\n", err)
		log.Fatal(err)
	}
	if err = command.Wait(); err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Printf("Error while reading. Error: %v\n", err)
		log.Fatal(err)
	}
	content := string(data)
	fmt.Printf("%v %d\n", content, len(content))
}

func testRender(cmd *cobra.Command, args []string) {
	snippet := snippet.SnippetData{
		ID:          1,
		Content:     "Test snippet",
		Description: "Test description",
		Tags:        "tag1,tag2",
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

func testPipe(cmd *cobra.Command, args []string) {
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

func testLanguages(cmd *cobra.Command, args []string) {
	fmt.Println("lexers:")
	sort.Sort(lexers.Registry.Lexers)
	for _, l := range lexers.Registry.Lexers {
		config := l.Config()
		fmt.Printf("  %s\n", config.Name)
		filenames := []string{}
		filenames = append(filenames, config.Filenames...)
		filenames = append(filenames, config.AliasFilenames...)
		if len(config.Aliases) > 0 {
			fmt.Printf("    aliases: %s\n", strings.Join(config.Aliases, " "))
		}
		if len(filenames) > 0 {
			fmt.Printf("    filenames: %s\n", strings.Join(filenames, " "))
		}
		if len(config.MimeTypes) > 0 {
			fmt.Printf("    mimetypes: %s\n", strings.Join(config.MimeTypes, " "))
		}
	}
	fmt.Println()
	fmt.Printf("styles:")
	for _, name := range styles.Names() {
		fmt.Printf(" %s", name)
	}
	fmt.Println()
	fmt.Printf("formatters:")
	for _, name := range formatters.Names() {
		fmt.Printf(" %s", name)
	}
	fmt.Println()
}

func testAdd(cmd *cobra.Command, args []string) {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)
	// s, err := testScan(au.Yellow("Snippet> ").String())
	// if err == testErrCanceled {
	// 	return
	// }
	// fmt.Println(s, err)
	description, err := testScan(au.Cyan("Description> ").String())
	if err == testErrCanceled {
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

func testScan(message string) (string, error) {
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
	return "", testErrCanceled
}

func init() {
	testCmd.AddCommand(testEditCmd)
	testCmd.AddCommand(testRenderCmd)
	testCmd.AddCommand(testPipeCmd)
	testCmd.AddCommand(testLanguagesCmd)
	testCmd.AddCommand(testAddCmd)
	rootCmd.AddCommand(testCmd)
}
