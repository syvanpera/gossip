package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/ui"
	"github.com/syvanpera/gossip/util"
)

var tags string
var stdin bool
var language string

var (
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a", "new", "n"},
		Short:   "Add a new bookmark",
		Long:    `Add a new bookmark"`,
		Args:    cobra.MaximumNArgs(2),
		Run:     add,
	}
)

func add(cmd *cobra.Command, args []string) {
	url := ""
	if len(args) > 0 {
		url = args[0]
	} else if url = ui.Prompt("URL", url); url == "" {
		return
	}

	if matched, _ := regexp.MatchString("^https?://*", url); !matched {
		url = "https://" + url
	}

	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	if description == "" || tags == "" {
		if meta := meta.Extract(url); meta != nil {
			if description == "" {
				description = meta.Description
			}
			if tags == "" {
				tags = fmt.Sprintf("%s,%s", tags, meta.Tags)
			}
		}
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return
		}
	}

	s, err := gossipService.Create(url, description, tags)
	if err != nil {
		util.PrintError(err)
		return
	}

	fmt.Println(s.Render(false))
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "tags")
}
