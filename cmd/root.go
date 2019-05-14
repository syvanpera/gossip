package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gossip",
	Short: "A command-line text snippet manager",
	Long: `gossip - A simple command-line text snippet manager.

Gossip is a CLI text snippet manager that can be used to store code blocks,
shell commands, bookmarks or any plain text snippets really.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gossip/config.toml)")
	rootCmd.PersistentFlags().BoolP("color", "c", false, "Toggle color output")
	viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.config/gossip")
	}

	viper.SetDefault("database", "example-db.json")
	viper.SetDefault("defaults.color", true)
	viper.SetDefault("types", []string{"code", "cmd", "url", "snippet"})

	viper.AutomaticEnv()
	viper.ReadInConfig()
}
