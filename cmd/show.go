package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get"},
	Short:   "Show a snippet",
	Long:    `Show a snippet`,
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

func show(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("Snippet with ID %d not found\n", id)
		return
	}

	fmt.Println(s)
}

func init() {
	rootCmd.AddCommand(showCmd)
}
