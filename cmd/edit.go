package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/snippet"
)

var addTags string

var editCmd = &cobra.Command{
	Use:     "edit ID",
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
	Run: editCommand,
}

func editCommand(_ *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	r := snippet.NewRepository()

	s := r.Get(id)
	if s == nil {
		fmt.Printf("Snippet with ID %d not found\n", id)
		return
	}

	existingTags := strings.Split(s.Data().Tags, ",")
	tags := strings.Split(addTags, ",")
	newTags := existingTags

	for _, t := range tags {
		tag := strings.ToLower(t)
		if !contains(existingTags, tag) {
			newTags = append(newTags, tag)
		}
	}

	s.Data().Tags = strings.Join(newTags, ",")
	r.Save(s)

	fmt.Println(s)
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&addTags, "tags", "t", "", "tags to add")
}
