package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var editTags string

var editCmd = &cobra.Command{
	Use:     "edit ID [CONTENT] [DESCRIPTION]",
	Aliases: []string{"ed", "e"},
	Short:   "Edit a bookmark",
	Long:    `Edit a bookmark"`,
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
	fmt.Println("ARGS", args, editTags)
	// id, _ := strconv.Atoi(args[0])

	// if editTags == "" {
	// 	util.PrintError(errors.New("Tags parameter missing"))
	// 	return
	// }

	// gossip, err := service.Update(id, editTags)
	// if err != nil {
	// 	util.PrintError(err)
	// 	return
	// }

	// fmt.Println(gossip.Render(false))
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&editTags, "tag", "t", "", "comma separated list of tags\n+xxx adds a tag, -xxx removes a tag")
}
