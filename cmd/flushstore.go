package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(flushstoreCMD)
}

var flushstoreCMD = &cobra.Command{
	Use:   "flush-store",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		UpdateStore(StoreConfig{StoreItems: []StoreItem{}})

	},
}
