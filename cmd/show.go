package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "s", "g"},
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

	s, err := service.GetSnippet(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(s)
}

func init() {
	rootCmd.AddCommand(showCmd)
}
