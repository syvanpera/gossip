package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/syvanpera/gossip/config"
	"github.com/syvanpera/gossip/internal/storage"
	"github.com/syvanpera/gossip/internal/ui"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open [id]",
	Short:   "Open a bookmark in your default web browser",
	Aliases: []string{"o"},
	Args:    cobra.ExactArgs(1), // Requires exactly one argument: the ID
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

		// 1. Fetch the bookmark
		bm, err := storage.GetByID(storagePath, id)
		if err != nil {
			fmt.Printf("%s %v\n", ui.StyleFailed.Render("✗ Error:"), err)
			os.Exit(1)
		}

		// 2. Print status
		fmt.Printf("%s Opening '%s'...\n", ui.StyleSuccess.Render("✓"), ui.StyleTitle.Render(bm.Title))

		// 3. Open the browser based on the Operating System
		err = openBrowser(bm.URL)
		if err != nil {
			fmt.Printf("%s Failed to open browser: %v\n", ui.StyleFailed.Render("✗"), err)
			os.Exit(1)
		}
	},
}

// openBrowser determines the OS and executes the appropriate command
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default: // Linux and standard Unix-like systems
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
