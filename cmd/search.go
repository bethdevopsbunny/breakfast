package cmd

import (
	"breakfast/store"
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

		storeConfig := store.RetrieveStoreConfig()

		for _, storeItem := range storeConfig.StoreItems {

			filepath := fmt.Sprintf("store/hash/%s/%s/%s", storeItem.Type, storeItem.Owner, storeItem.Repo)

			// for each directory
			// for each directory provided

			hashDirectories, _ := ioutil.ReadDir(filepath)

			for _, dir := range hashDirectories {

				filepaths := fmt.Sprintf("%s/%s", filepath, dir.Name())

				files, _ := ioutil.ReadDir(filepaths)

				for _, i := range files {

					filel := fmt.Sprintf("%s/%s", filepaths, i.Name())

					item := retriveHashFile(filel)

					for _, j := range item {
						if j.Hash == args[0] {
							log.Infof("File - %s", i.Name())
							println(j.Pass)

						}

					}

				}
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

func listHashStoreDirs(items []store.StoreItem) {

}
