package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"hcw/hashes"
	github "hcw/source/github"
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

				data, _ := github.GetLatestReleaseData(element.Owner, element.Repo)
				published := data.PublishedAt.UnixMicro()

				createStoreItemHash(element)
				store := retrieveStoreConfig()
				addToStore(store, element, published, createStoreItemHash(element))

				println(isItInStore(store, element))

				//
				//zipfilepath := fmt.Sprintf("store/zip/%s-%s-%s-%d.zip", element.Type, element.Owner, element.Repo, published)
				//
				//if _, err := os.Stat(zipfilepath); err == nil {
				//
				//	println("Already Uptodate")
				//
				//} else {
				//
				//	github.DownloadFile(zipfilepath, data.ZipballURL)
				//
				//	//
				//	//txtfilepath := fmt.Sprintf("store/txt/%d-words", zipfilepath)
				//	//zipRoot := source.Unzip(zipfilepath, "store/txt")
				//	//zipRootPath := fmt.Sprintf("store/txt/%s", zipRoot)
				//	//os.Rename(zipRootPath, txtfilepath)
				//	//println("completed")
				//
				//}

			}

		}

		//for _, element := range config.Wordlistrepos[0].Includedfiles {
		//	sa := ReadEachLine(element)
		//	for _, element2 := range sa {
		//		println(fmt.Sprintf("%s:%s", element2, hashes.SHA256(element2)))
		//	}
		//
		//}

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
	Dateadded int64  `json:"dateadded"`
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

		if element.Hash == createStoreItemHash(wordlist) {

			isIt = true
		}

	}
	return isIt

}

func createStoreItemHash(wordlist Wordlist) string {
	wordliststring := fmt.Sprintf("%s%s%s", wordlist.Type, wordlist.Owner, wordlist.Repo)
	return hashes.SHA1(wordliststring)

}

func addToStore(storeConfig StoreConfig, wordlist Wordlist, published int64, storeItemHash string) {

	storeConfig.StoreItems = append(storeConfig.StoreItems, StoreItem{
		Type:      wordlist.Type,
		Owner:     wordlist.Owner,
		Repo:      wordlist.Repo,
		Dateadded: published,
		Filename:  "nothing yet",
		Hash:      storeItemHash,
	})

	file, _ := json.MarshalIndent(storeConfig, "", " ")

	_ = ioutil.WriteFile("store/storeconfig.json", file, 0644)

}
