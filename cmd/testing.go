package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(testingCmd)
}

var testingCmd = &cobra.Command{}
