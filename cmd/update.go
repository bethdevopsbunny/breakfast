package cmd

import (
	"breakfast/hashes"
	"breakfast/source"
	"breakfast/source/github"
	"breakfast/store"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func init() {

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{

	Use:   "update",
	Short: "checks for and updates store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		config := retrieveConfig()

		for _, sourcePack := range config.SourcePacks {

			if sourcePack.Type == "github" {

				data, err := github.GetLatestReleaseData(sourcePack.Owner, sourcePack.Repo)
				if err != nil {
					log.Errorf("Failed to connect to github %s repo %s", sourcePack.Owner, sourcePack.Repo)
				}
				published := data.PublishedAt.UnixMicro()
				storeConfig := store.RetrieveStoreConfig()

				storeCheck := false

				for i, storeItem := range storeConfig.StoreItems {

					if store.IsItInStore(storeItem, storeItem.Type, storeItem.Owner, storeItem.Repo) {

						storeCheck = true

						if !areSourcePackWordListsUptodate(storeItem, sourcePack) {

							updateSourcePackWordLists(storeConfig, i, sourcePack)

							for _, desiredHash := range config.Global.EncryptionHashes {

								wordListHashUpdate(desiredHash, storeConfig)
							}

							// run wordlist hash update

						}
					}
				}

				if !storeCheck {

					fullHash := fmt.Sprintf("%x", store.Hash(sourcePack))

					sourceHash := fmt.Sprintf("%x", store.Hash(sourcePack.Type, sourcePack.Owner, sourcePack.Repo))
					includedFilesHash := fmt.Sprintf("%x", store.Hash(sourcePack.IncludedWordLists))

					store.AddToStore(storeConfig, sourcePack, published, store.StoreHashes{fullHash, sourceHash, includedFilesHash})

					updateSourcePacks(sourcePack, data)

				}

			}

		}

	},
}

func updateSourcePacks(pack source.SourcePack, data github.ReleaseData) {

	filename := fmt.Sprintf("%s-%s-%s", pack.Type, pack.Owner, pack.Repo)
	zipfilepath := fmt.Sprintf("store/zip/%s.zip", filename)
	//github.DownloadFile(zipfilepath, data.ZipballURL)
	txtfilepath := fmt.Sprintf("store/txt/%s", filename)
	zipRoot := source.Unzip(zipfilepath, "store/txt")
	zipRootPath := fmt.Sprintf("store/txt/%s", zipRoot)
	os.Rename(zipRootPath, txtfilepath)
	log.Infof("Complted")

}

type HashedPasswordItem struct {
	Pass string `json:"pass"`
	Hash string `json:"hash"`
}

func wordListHashUpdate(hashtype string, storeConfig store.StoreConfig) {

	// four for loops sucks. remove this then your remove them.

	for _, storeItem := range storeConfig.StoreItems {

		for _, wordList := range storeItem.IncludedWordLists {

			if fileExists(wordList) {
				var hashedPasswordItems []HashedPasswordItem
				passwords := ReadEachLine(wordList)

				for _, password := range passwords {

					switch hashtype {
					case "SHA1":
						hashedPasswordItems = append(hashedPasswordItems, HashedPasswordItem{Pass: password, Hash: hashes.SHA1(password)})
						log.Infof("SHA1 Hashing %s", password)
					case "SHA256":
						hashedPasswordItems = append(hashedPasswordItems, HashedPasswordItem{Pass: password, Hash: hashes.SHA256(password)})
						log.Infof("SHA256 Hashing %s", password)
					case "MD5":
						hashedPasswordItems = append(hashedPasswordItems, HashedPasswordItem{Pass: password, Hash: hashes.MD5Hash(password)})
						log.Infof("MD5 Hashing %s", password)
					case "BCRYPT":
						hashedPasswordItems = append(hashedPasswordItems, HashedPasswordItem{Pass: password, Hash: hashes.BCRYPT(password)})
						log.Infof("BCRYPT Hashing %s", password)
					}
				}
				UpdateList(hashedPasswordItems, storeItem, wordList, hashtype)
			}
		}
	}

}

func fileExists(filePath string) bool {

	if _, err := os.Stat(filePath); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		log.Infof("unable to find file - %s", filePath)
		return false
	} else {
		log.Errorf("corruption with file - %s", filePath)
		return false
	}
}

func updateSourcePackWordLists(storeConfig store.StoreConfig, index int, sourcePack source.SourcePack) {

	storeConfig.StoreItems[index].IncludedWordLists = sourcePack.IncludedWordLists
	storeConfig.StoreItems[index].StoreHashes.IncludedWordListItems = fmt.Sprintf("%x", store.Hash(sourcePack.IncludedWordLists))
	store.UpdateStore(storeConfig)

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
	SourcePacks []source.SourcePack `json:"sourcepacks"`
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

func areSourcePackWordListsUptodate(storeItem store.StoreItem, sourcePack source.SourcePack) (isIt bool) {

	if storeItem.StoreHashes.IncludedWordListItems == fmt.Sprintf("%x", store.Hash(sourcePack.IncludedWordLists)) {

		println("Hashes match")
		isIt = true

	}

	return isIt

}

func UpdateList(out []HashedPasswordItem, storeItem store.StoreItem, wordlistFilename string, hash string) {

	filepath := fmt.Sprintf("store/hash/%s/%s/%s/%s", storeItem.Type, storeItem.Owner, storeItem.Repo, hash)
	err := os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	filename := fmt.Sprintf("%s/%s.json", filepath, filenameFromFilepath(wordlistFilename))

	file, _ := json.MarshalIndent(out, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

}

func filenameFromFilepath(filepath string) string {

	split := strings.Split(filepath, "/")
	single := split[len(split)-1]

	singleSplit := strings.Split(single, ".")

	return singleSplit[0]
}
