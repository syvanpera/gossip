// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}

	// Set up the configuration path
	configPath := filepath.Join(home, ".config", "gossip")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("storage_file", filepath.Join(configPath, "bookmarks.json"))

	// Read or create the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a default one
			viper.SafeWriteConfig()
		} else {
			fmt.Println("Error reading config:", err)
		}
	}
}

// GetStoragePath is a helper to easily grab the JSON file location
func GetStoragePath() string {
	return viper.GetString("storage_file")
}
