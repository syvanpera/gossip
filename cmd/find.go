package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"f", "get", "g", "show", "s"},
	Short:   "Find a snippet",
	Long:    `Find a snippet`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: find,
}

func find(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])

	gossip, err := gossipService.Get(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(gossip.Render())
}

func init() {
	rootCmd.AddCommand(findCmd)
}
