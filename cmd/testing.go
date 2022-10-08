package cmd

import (
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

		arr := [7]string{"This", "is", "the", "tutorial",
			"of", "Go", "language"}

		fmt.Printf("%x\n", Hash(arr))
		fmt.Printf("%x\n", Hash("123"))

	},
}

//func Hash(objs ...interface{}) []byte {
//	digester := crypto.SHA1.New()
//	for _, ob := range objs {
//		fmt.Fprint(digester, reflect.TypeOf(ob))
//		fmt.Fprint(digester, ob)
//	}
//	return digester.Sum(nil)
//}
