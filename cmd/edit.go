// cmd/edit.go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/fetcher"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var (
	editTitle   string
	editURL     string
	editComment string
	editTags    []string
	refetchTags bool
	refetchAll  bool
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit an existing bookmark's title, url, tags, or comment",
	Args:  cobra.ExactArgs(1),
	// Dynamic TAB completion for the ID
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		storagePath := config.GetStoragePath()
		bookmarks, err := storage.Load(storagePath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, b := range bookmarks {
			if strings.HasPrefix(b.ID, toComplete) {
				completion := fmt.Sprintf("%s\t%s", b.ID, b.Title)
				completions = append(completions, completion)
			}
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		storagePath := config.GetStoragePath()

		// 1. Fetch the existing bookmark
		bm, err := storage.GetByID(storagePath, id)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✗ Error:"), err)
			os.Exit(1)
		}

		// 2. Handle Refetching FIRST
		if refetchTags || refetchAll {
			fmt.Printf("\n%s Fetching metadata for '%s'...\n", ui.StyleRunning.Render("Working..."), ui.StyleURL.Render(bm.URL))

			f := fetcher.GetFetcher(bm.URL)
			meta, err := f.Fetch(bm.URL)
			if err != nil {
				fmt.Printf("\n%s Failed to fetch metadata: %v\n", ui.StyleFailed.Render("✗"), err)
				os.Exit(1)
			}

			if refetchAll {
				if meta.Title != "" {
					bm.Title = meta.Title
				}
				if meta.Description != "" {
					bm.Comment = meta.Description
				}
				if len(meta.Tags) > 0 {
					bm.Tags = meta.Tags
				}
				fmt.Printf("%s Refetched all metadata successfully!\n", ui.StyleSuccess.Render("✓"))
			} else if refetchTags {
				if len(meta.Tags) > 0 {
					bm.Tags = meta.Tags
					fmt.Printf("%s Refetched tags successfully!\n", ui.StyleSuccess.Render("✓"))
				} else {
					fmt.Printf("\n%s No default tags found on the target URL.\n", ui.StyleComment.Render("-"))
				}
			}
		}

		// 3. Apply manual updates (These safely override the refetched data)
		if cmd.Flags().Changed("title") {
			bm.Title = editTitle
		}
		if cmd.Flags().Changed("url") {
			bm.URL = editURL
		}
		if cmd.Flags().Changed("comment") {
			bm.Comment = editComment
		}
		if cmd.Flags().Changed("tags") {
			bm.Tags = editTags
		}

		// 4. Save the updated bookmark back to storage
		err = storage.Update(storagePath, bm)
		if err != nil {
			fmt.Printf("%s Failed to update: %v\n", ui.StyleFailed.Render("✖"), err)
			os.Exit(1)
		}

		fmt.Printf("%s Bookmark '%s' updated successfully!\n", ui.StyleSuccess.Render("✓ Success!"), ui.StyleID.Render(id))
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&editTitle, "title", "T", "", "Update the title")
	editCmd.Flags().StringVarP(&editURL, "url", "u", "", "Update the URL")
	editCmd.Flags().StringVarP(&editComment, "comment", "c", "", "Update the comment")
	editCmd.Flags().StringSliceVarP(&editTags, "tags", "t", []string{}, "Update tags (comma-separated)")

	editCmd.Flags().BoolVar(&refetchTags, "refetch-tags", false, "Refetch default tags from the URL")
	editCmd.Flags().BoolVar(&refetchAll, "refetch-all", false, "Refetch all metadata (title, comment, tags) from the URL") // Register new flag
}
