package snippet

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

type Command struct {
	data SnippetData
}

func (cmd *Command) Type() SnippetType  { return COMMAND }
func (cmd *Command) Data() *SnippetData { return &cmd.data }

func (c *Command) Execute() error {
	fmt.Print(c)

	prompt := promptui.Prompt{
		Label:     "Are you sure you want to execute this command",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		fmt.Println("Canceled")
		return err
	}

	var command *exec.Cmd
	command = exec.Command("sh", "-c", c.data.Content)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	command.Run()

	return nil
}

func (c *Command) String() string {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder

	fmt.Fprintf(&output, "\n%s ", au.BrightCyan(fmt.Sprintf("%d.", c.data.ID)))
	fmt.Fprintln(&output, au.Bold(au.BrightGreen(c.data.Description)))
	fmt.Fprintf(&output, "   %s %s\n", au.BrightRed(">"), au.BrightYellow(c.data.Content))
	if c.data.Tags != "" {
		fmt.Fprintf(&output, "   %s %s\n", au.BrightRed("#"), au.BrightBlue(c.data.Tags))
	}

	return output.String()
}

func NewCommand(command, description string, tags string) *Command {
	cmd := Command{
		data: SnippetData{
			Content:     command,
			Description: description,
			Tags:        tags,
			Type:        COMMAND,
		},
	}

	return &cmd
}
