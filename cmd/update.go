package cmd

import (
	"breakfast/hashes"
	"breakfast/source"
	github "breakfast/source/github"
	"bufio"
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
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

		for _, sourcePack := range config.SourcePacks {

			if sourcePack.Type == "github" {

				data, _ := github.GetLatestReleaseData(sourcePack.Owner, sourcePack.Repo)
				published := data.PublishedAt.UnixMicro()
				store := retrieveStoreConfig()

				storeCheck := false

				for i, storeItem := range store.StoreItems {

					if isItInStore(storeItem, storeItem.Type, storeItem.Owner, storeItem.Repo) {

						storeCheck = true

						if !areSourcePackWordListsUptodate(storeItem, sourcePack) {

							updateSourcePackWordLists(store, i, sourcePack)
							wordListHashUpdate(config, store)
							// run wordlist hash update

						}
					}
				}

				if !storeCheck {
					fullHash := fmt.Sprintf("%x", Hash(sourcePack))

					sourceHash := fmt.Sprintf("%x", Hash(sourcePack.Type, sourcePack.Owner, sourcePack.Repo))
					includedFilesHash := fmt.Sprintf("%x", Hash(sourcePack.IncludedWordLists))

					addToStore(store, sourcePack, published, StoreHashes{fullHash, sourceHash, includedFilesHash})

					updateSourcePacks(sourcePack, data)

				}

				// do not run hash update

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
				//
				//

				//}

			}

		}

		//for _, element := range config.Wordlistrepos[0].includedWordLists {
		//	sa := ReadEachLine(element)
		//	for _, element2 := range sa {
		//		println(fmt.Sprintf("%s:%s", element2, hashes.SHA256(element2)))
		//	}
		//
		//}

	},
}

func updateSourcePacks(pack SourcePack, data github.ReleaseData) {

	filename := fmt.Sprintf("%s-%s-%s", pack.Type, pack.Owner, pack.Repo)
	zipfilepath := fmt.Sprintf("store/zip/%s.zip", filename)
	//github.DownloadFile(zipfilepath, data.ZipballURL)
	txtfilepath := fmt.Sprintf("store/txt/%s", filename)
	zipRoot := source.Unzip(zipfilepath, "store/txt")
	zipRootPath := fmt.Sprintf("store/txt/%s", zipRoot)
	os.Rename(zipRootPath, txtfilepath)
	println("completed")

}

type HashedPasswordList struct {
	Pass string `json:"pass"`
	Hash string `json:"hash"`
}

func wordListHashUpdate(config Config, storeConfig StoreConfig) {

	for _, e := range config.Global.EncryptionHashes {

		if e == "SHA1" {

			var o []HashedPasswordList

			for _, i := range storeConfig.StoreItems {

				for _, j := range i.IncludedWordLists {

					sa := ReadEachLine(j)

					for _, k := range sa {
						o = append(o, HashedPasswordList{Pass: k, Hash: hashes.SHA1(k)})
						//println(fmt.Sprintf("%s:%s", k, hashes.SHA1(k)))
					}

					UpdateList(o, i, j, "SHA1")

				}

			}

		}

		if e == "SHA256" {

			var o []HashedPasswordList

			for _, i := range storeConfig.StoreItems {

				for _, j := range i.IncludedWordLists {

					sa := ReadEachLine(j)

					for _, k := range sa {
						o = append(o, HashedPasswordList{Pass: k, Hash: hashes.SHA256(k)})
						//println(fmt.Sprintf("%s:%s", k, hashes.SHA1(k)))
					}
					UpdateList(o, i, j, "SHA256")
				}

			}

		}
		if e == "MD5" {

			var o []HashedPasswordList

			for _, i := range storeConfig.StoreItems {

				for _, j := range i.IncludedWordLists {

					sa := ReadEachLine(j)

					for _, k := range sa {
						o = append(o, HashedPasswordList{Pass: k, Hash: hashes.MD5Hash(k)})
						//println(fmt.Sprintf("%s:%s", k, hashes.SHA1(k)))
					}
					UpdateList(o, i, j, "MD5")
				}

			}

		}

	}

}

func updateSourcePackWordLists(store StoreConfig, index int, sourcePack SourcePack) {

	store.StoreItems[index].IncludedWordLists = sourcePack.IncludedWordLists
	store.StoreItems[index].StoreHashes.IncludedWordLists = fmt.Sprintf("%x", Hash(sourcePack.IncludedWordLists))
	UpdateStore(store)

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
	Global struct {
		EncryptionHashes []string `json:"encryptionhashes"`
	} `json:"global"`
	SourcePacks []SourcePack `json:"sourcepacks"`
}

type SourcePack struct {
	Type              string   `json:"type"`
	Owner             string   `json:"owner"`
	Repo              string   `json:"repo"`
	IncludedWordLists []string `json:"includedwordlists"`
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
	Type              string      `json:"type"`
	Owner             string      `json:"owner"`
	Repo              string      `json:"repo"`
	Dateadded         int64       `json:"dateadded"`
	IncludedWordLists []string    `json:"includedWordLists"`
	StoreHashes       StoreHashes `json:"storeHashes"`
}

type StoreHashes struct {
	Full              string `json:"full"`
	Source            string `json:"source"`
	IncludedWordLists string `json:"includedwordlists"`
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

func areSourcePackWordListsUptodate(storeItem StoreItem, sourcePack SourcePack) (isIt bool) {

	if storeItem.StoreHashes.IncludedWordLists == fmt.Sprintf("%x", Hash(sourcePack.IncludedWordLists)) {

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

func UpdateStore(newStore StoreConfig) {

	file, _ := json.MarshalIndent(newStore, "", " ")
	_ = ioutil.WriteFile("store/store.json", file, 0644)

}

func UpdateList(out []HashedPasswordList, storeItem StoreItem, wordlistFilename string, hash string) {

	filename := fmt.Sprintf("store/hash/%s-%s-%s-%s-%s.json", storeItem.Type, storeItem.Owner, storeItem.Repo, filenameFromFilepath(wordlistFilename), hash)

	file, _ := json.MarshalIndent(out, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

}

func addToStore(storeConfig StoreConfig, sourcePack SourcePack, published int64, hashes StoreHashes) {

	storeConfig.StoreItems = append(storeConfig.StoreItems, StoreItem{
		Type:              sourcePack.Type,
		Owner:             sourcePack.Owner,
		Repo:              sourcePack.Repo,
		Dateadded:         published,
		IncludedWordLists: sourcePack.IncludedWordLists,
		StoreHashes:       StoreHashes{Full: hashes.Full, Source: hashes.Source, IncludedWordLists: hashes.IncludedWordLists},
	})

	file, _ := json.MarshalIndent(storeConfig, "", " ")

	_ = ioutil.WriteFile("store/store.json", file, 0644)

}

func filenameFromFilepath(filepath string) string {

	split := strings.Split(filepath, "/")
	single := split[len(split)-1]

	singleSplit := strings.Split(single, ".")

	return singleSplit[0]
}
