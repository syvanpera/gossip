package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"
)

var appName = "gossip"

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "A command-line text snippet manager",
	Long: `gossip - A simple command-line text snippet manager.

Gossip is a CLI text snippet manager that can be used to store code blocks,
shell commands, bookmarks or any plain text snippets really.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("color", "c", false, "toggle color output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s/%s", configPath(), appName))

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("config.color", true)
	viper.SetDefault("config.browser", "default")
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Unable to resolve user config directory.", err)
	}

	return dir
}

func dataPath() string {
	dir, err := util.UserDataDir()
	if err != nil {
		log.Fatal("Unable to resolve user data directory.", err)
	}

	return dir
}
