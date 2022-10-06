package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"hcw/hashes"
	intfgeneral "hcw/intf/general"
	intfgithub "hcw/intf/github"
	"io/ioutil"
	"os"
)

func init() {

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{

	Use:   "update",
	Short: "checks for and updates store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// collect config
		// check for changes to source
		// pull down and prepare each new source.
		// for each file stipulated in config hash with chosen hashes

		config := retrieveConfig()

		for _, element := range config.Wordlists {

			if element.Type == "github" {

				data, _ := intfgithub.GetLatestReleaseData(element.Owner, element.Repo)
				published := data.PublishedAt.UnixMicro()

				s := retrieveStoreConfig()
				isItInStore(s, element)

				zipfilepath := fmt.Sprintf("store/zip/%s-%s-%s-%d.zip", element.Type, element.Owner, element.Repo, published)

				if _, err := os.Stat(zipfilepath); err == nil {

					println("Already Uptodate")

				} else {

					intfgithub.DownloadFile(zipfilepath, data.ZipballURL)

					txtfilepath := fmt.Sprintf("store/txt/%d-words", zipfilepath)
					zipRoot := intfgeneral.Unzip(zipfilepath, "store/txt")
					zipRootPath := fmt.Sprintf("store/txt/%s", zipRoot)
					os.Rename(zipRootPath, txtfilepath)
					println("completed")

				}

			}

		}

		for _, element := range config.Wordlistrepos[0].Includedfiles {
			sa := ReadEachLine(element)
			for _, element2 := range sa {
				println(fmt.Sprintf("%s:%s", element2, hashes.SHA256(element2)))
			}

		}

	},
}

func ReadEachLine(filepath string) (fileLines []string) {

	readFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

type Config struct {
	Wordlists []Wordlist `json:"wordlists"`
}

type Wordlist struct {
	Type          string   `json:"type"`
	Owner         string   `json:"owner"`
	Repo          string   `json:"repo"`
	Includedfiles []string `json:"includedfiles"`
}

func retrieveConfig() Config {

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	c, err := ioutil.ReadAll(jsonFile)
	if err != nil {

	}

	var config Config

	json.Unmarshal(c, &config)

	return config

}

type StoreConfig struct {
	StoreItems []StoreItem `json:"StoreItems"`
}

type StoreItem struct {
	Type      string `json:"type"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Dateadded string `json:"dateadded"`
	Filename  string `json:"filename"`
	Hash      string `json:"hash"`
}

func retrieveStoreConfig() StoreConfig {

	jsonFile, err := os.Open("store/storeconfig.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	c, err := ioutil.ReadAll(jsonFile)
	if err != nil {

	}

	var storeConfig StoreConfig

	json.Unmarshal(c, &storeConfig)

	return storeConfig

}

func isItInStore(storeConfig StoreConfig, wordlist Wordlist) (isIt bool) {

	for _, element := range storeConfig.StoreItems {

		if element.Hash == wordlist.hash {

		}

	}

}
