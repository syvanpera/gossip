package snippet

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/ui"
)

type Command struct {
	data SnippetData
}

func (cmd *Command) Type() SnippetType  { return COMMAND }
func (cmd *Command) Data() *SnippetData { return &cmd.data }

func (c *Command) Execute() error {
	fmt.Println(c)

	if !ui.Confirm("Are you sure you want to execute this command") {
		return ErrExecCanceled
	}

	var command *exec.Cmd
	command = exec.Command("sh", "-c", c.data.Content)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	command.Run()

	return nil
}

func (c *Command) Edit() error {
	s, err := Edit("Command", c.Data().Content)
	if err != nil {
		return err
	}
	c.Data().Content = s

	s, err = Edit("Description", c.Data().Description)
	if err != nil {
		return err
	}
	c.Data().Description = s

	return nil
}

func (c *Command) String() string {
	colors := viper.GetBool("defaults.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	var output strings.Builder

	fmt.Fprintf(&output, "\n%s ", au.BrightCyan(fmt.Sprintf("%d.", c.data.ID)))
	fmt.Fprintln(&output, au.Bold(au.BrightGreen(c.data.Description)))
	fmt.Fprintf(&output, "   %s %s", au.BrightRed(">"), au.BrightYellow(c.data.Content))
	if c.data.Tags != "" {
		fmt.Fprintf(&output, "\n   %s %s", au.BrightRed("#"), au.BrightBlue(c.data.Tags))
	}

	return output.String()
}
