package cmd

import (
	"breakfast/hashes"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(testingCmd)
}

var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Println(hashes.BCRYPT("hello"))

		fmt.Println(hashes.BCRYPTTEST("$2a$10$FwMqin3GEaOf8xcEduQPfu4.l3F0jT00DLOvH9SK89kbwWBRTePSC", "helo"))

	},
}
