// cmd/root.go
package cmd

import (
	"os"

	"github.com/syvanpera/gossip/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gossip",
	Short: "A beautiful command-line bookmark manager",
	Long:  `A lightweight, terminal-based bookmark manager built with Go, Cobra, and Bubble Tea.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration when the app starts
	cobra.OnInitialize(config.InitConfig)
}
