package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hcw",
	Short: "Hash Crack Worker",
	Long:  " ",
}

func init() {

}

func Execute() error {
	return rootCmd.Execute()
}
