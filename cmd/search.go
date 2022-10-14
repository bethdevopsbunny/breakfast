package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func init() {

	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		item := retriveHashFile()

		for _, i := range item {
			if i.Hash == args[0] {
				println(i.Pass)
			}

		}

	},
}

func retriveHashFile() []HashedPasswordItem {

	jsonFile, err := os.Open("store/hash/github-danielmiessler-SecLists-500-worst-passwords-MD5.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	c, err := ioutil.ReadAll(jsonFile)
	if err != nil {

	}

	var hashedPasswordItem []HashedPasswordItem

	json.Unmarshal(c, &hashedPasswordItem)

	return hashedPasswordItem

}
