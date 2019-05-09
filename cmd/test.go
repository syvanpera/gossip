package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	Long:  `Testing`,
	RunE:  test,
}

func test(cmd *cobra.Command, args []string) error {
	fmt.Printf("test: %v\n", args)

	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
