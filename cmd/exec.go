package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var execCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"run"},
	Short:   "Execute a command snippet",
	Long:    `Execute a command snippet`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: execute,
}

func execute(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("snippet #%d not found\n", id)
		return
	}

	s.Execute()
}

func init() {
	rootCmd.AddCommand(execCmd)
}
