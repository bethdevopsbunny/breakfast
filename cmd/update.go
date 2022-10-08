package cmd

import (
	"bufio"
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	github "hcw/source/github"
	"io/ioutil"
	"os"
	"reflect"
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

		for _, wordlist := range config.Wordlists {

			if wordlist.Type == "github" {

				data, _ := github.GetLatestReleaseData(wordlist.Owner, wordlist.Repo)
				published := data.PublishedAt.UnixMicro()
				store := retrieveStoreConfig()

				storeCheck := false
				for i, storeItem := range store.StoreItems {

					if isItInStore(storeItem, storeItem.Type, storeItem.Owner, storeItem.Repo) {

						storeCheck = true

						if !isWordListsUpToDate(storeItem, wordlist) {

							store.StoreItems[i].Includedfiles = wordlist.Includedfiles
							updateStore(store)
						}

					}

				}

				if !storeCheck {
					fullHash := fmt.Sprintf("%x", Hash(wordlist))
					sourceHash := fmt.Sprintf("%x", Hash(wordlist.Type, wordlist.Owner, wordlist.Repo))
					wordlistHash := fmt.Sprintf("%x", Hash(wordlist.Includedfiles))

					addToStore(store, wordlist, published, StoreHashes{fullHash, sourceHash, wordlistHash})

				}

				// does the source exist

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
	Type          string      `json:"type"`
	Owner         string      `json:"owner"`
	Repo          string      `json:"repo"`
	Dateadded     int64       `json:"dateadded"`
	Includedfiles []string    `json:"includedfiles"`
	StoreHashes   StoreHashes `json:"storeHashes"`
}

type StoreHashes struct {
	Full     string `json:"full"`
	Source   string `json:"source"`
	WordList string `json:"wordlist"`
}

func retrieveStoreConfig() StoreConfig {

	jsonFile, err := os.Open("store/store.json")
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

func isItInStore(storeItem StoreItem, Type string, Owner string, Repo string) (isIt bool) {

	if storeItem.StoreHashes.Source == fmt.Sprintf("%x", Hash(Type, Owner, Repo)) {
		isIt = true
	}

	return isIt

}

func isWordListsUpToDate(storeItem StoreItem, wordlist Wordlist) (isIt bool) {

	if storeItem.StoreHashes.WordList == fmt.Sprintf("%x", Hash(wordlist.Includedfiles)) {

		println("Hashes match")
		isIt = true

	}

	return isIt

}

func Hash(objs ...interface{}) []byte {
	digester := crypto.SHA1.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return digester.Sum(nil)
}

func updateStore(newStore StoreConfig) {

	file, _ := json.MarshalIndent(newStore, "", " ")
	_ = ioutil.WriteFile("store/store.json", file, 0644)

}

func addToStore(storeConfig StoreConfig, wordlist Wordlist, published int64, hashes StoreHashes) {

	storeConfig.StoreItems = append(storeConfig.StoreItems, StoreItem{
		Type:          wordlist.Type,
		Owner:         wordlist.Owner,
		Repo:          wordlist.Repo,
		Dateadded:     published,
		Includedfiles: wordlist.Includedfiles,
		StoreHashes:   StoreHashes{Full: hashes.Full, Source: hashes.Source, WordList: hashes.WordList},
	})

	file, _ := json.MarshalIndent(storeConfig, "", " ")

	_ = ioutil.WriteFile("store/store.json", file, 0644)

}
