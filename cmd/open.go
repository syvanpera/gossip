package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open a bookmark",
	Long:    `Open a bookmark`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: open,
}

func open(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])

	gossipService.Open(id)
}

func init() {
	rootCmd.AddCommand(openCmd)
}
