package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/pkg/db"
	"github.com/syvanpera/gossip/pkg/services/bookmarks"
	"github.com/syvanpera/gossip/pkg/services/bookmarks/store"
)

var service bookmarks.Service

var rootCmd = &cobra.Command{
	Use:   "gossip",
	Short: "A command-line bookmark manager",
	Long:  `gossip - A simple command-line bookmark manager.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	conn := db.GetConnection(viper.GetString("database.path"))
	defer conn.Close()

	service = bookmarks.New(store.New(conn))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Command failed")
	}
}

func init() {
	// rootCmd.PersistentFlags().BoolP("color", "c", false, "turn off color output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "turn on debug mode")

	// viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
