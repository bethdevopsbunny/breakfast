package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		testFile := "store/hash/github-danielmiessler-SecLists-500-worst-passwords-MD5.json"
		item := retriveHashFile(testFile)

		for _, i := range item {
			if i.Hash == args[0] {
				log.Infof("File - %s", testFile)
				println(i.Pass)

			}

		}

	},
}

func retriveHashFile(filepath string) []HashedPasswordItem {

	jsonFile, err := os.Open(filepath)
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
