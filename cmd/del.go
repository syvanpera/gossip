package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var force bool

var delCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"rm"},
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
		fmt.Printf("Snippet with ID %d not found\n", id)
		return
	}
	fmt.Println("Got", s)

	// if err := r.Del(id); err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("Snippet removed...")
}

func init() {
	delCmd.Flags().BoolVarP(&force, "force", "f", false, `Don't ask for confirmation`)
	rootCmd.AddCommand(delCmd)
}
