package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/gossip"
	"github.com/syvanpera/gossip/snippet"
	"github.com/syvanpera/gossip/util"
)

var appName = "gossip"

var service snippet.Service
var gossipService gossip.Service

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
	cobra.OnInitialize(initServices)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Clean(fmt.Sprintf("%s/%s", configPath(), appName)))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("config.color", true)
	viper.SetDefault("config.browser", "default")
	viper.SetDefault("config.database", filepath.Clean(fmt.Sprintf("%s/%s/%s.db", dataPath(), appName, appName)))

	rootCmd.PersistentFlags().BoolP("color", "c", false, "toggle color output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "be quiet")

	viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initServices() {
	service = snippet.NewService(snippet.NewSQLiteRepository(viper.GetString("config.database")))
	gossipService = gossip.NewService(viper.GetString("config.database"))
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
