package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/util"
)

var addTags string

var editCmd = &cobra.Command{
	Use:     "edit ID [CONTENT] [DESCRIPTION]",
	Aliases: []string{"ed", "e"},
	Short:   "Edit a snippet",
	Long:    `Edit a snippet"`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an integer ID argument")
		}
		if _, err := strconv.Atoi(args[0]); err == nil {
			return nil
		}
		return fmt.Errorf("invalid ID number specified: %s", args[0])
	},
	Run: edit,
}

func edit(_ *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])

	if addTags == "" {
		util.PrintError(errors.New("Tags parameter missing"))
		return
	}

	gossip, err := gossipService.Update(id, addTags)
	if err != nil {
		util.PrintError(err)
		return
	}

	fmt.Println(gossip.Render())
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&addTags, "tags", "t", "", "tags to add")
}
