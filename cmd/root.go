package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "breakfast",
	Short: "dynamically updated store of popular passwords hashed with popular hashes",
	Long:  " ",
}

func init() {

}

func Execute() error {
	return rootCmd.Execute()
}
