package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/ui"
)

var force bool

var delCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"rm", "d"},
	Short:   "Delete a snippet",
	Long:    `Delete a snippet`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: del,
}

func del(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("Snippet #%d not found\n", id)
		return
	}

	fmt.Println(s)

	if !force && !ui.Confirm("Are you sure you want to delete this snippet") {
		fmt.Println("Canceled")
		return
	}

	r.Del(id)

	fmt.Println("Snippet deleted...")
}

func init() {
	delCmd.Flags().BoolVarP(&force, "force", "f", false, `Don't ask for confirmation`)
	rootCmd.AddCommand(delCmd)
}
