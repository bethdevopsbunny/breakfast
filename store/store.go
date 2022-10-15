package store

import (
	"breakfast/source"
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type StoreConfig struct {
	StoreItems []StoreItem `json:"StoreItems"`
}

type StoreItem struct {
	Type              string      `json:"type"`
	Owner             string      `json:"owner"`
	Repo              string      `json:"repo"`
	Dateadded         int64       `json:"dateadded"`
	IncludedWordLists []string    `json:"includedWordListItems"`
	StoreHashes       StoreHashes `json:"storeHashes"`
}

type StoreHashes struct {
	Full                  string `json:"full"`
	Source                string `json:"source"`
	IncludedWordListItems string `json:"includedwordlistitems"`
}

type IncludedWordListItem struct {
	NonHashTextFilepath string `json:"nonhashtextfilepath"`
	StoreLocation       string `json:"storeLocation"`
}

func RetrieveStoreConfig() StoreConfig {

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

func IsItInStore(storeItem StoreItem, Type string, Owner string, Repo string) (isIt bool) {

	if storeItem.StoreHashes.Source == fmt.Sprintf("%x", Hash(Type, Owner, Repo)) {
		isIt = true
	}

	return isIt

}

func UpdateStore(newStore StoreConfig) {

	file, _ := json.MarshalIndent(newStore, "", " ")
	_ = ioutil.WriteFile("store/store.json", file, 0644)

}

func AddToStore(storeConfig StoreConfig, sourcePack source.SourcePack, published int64, hashes StoreHashes) {

	storeConfig.StoreItems = append(storeConfig.StoreItems, StoreItem{
		Type:              sourcePack.Type,
		Owner:             sourcePack.Owner,
		Repo:              sourcePack.Repo,
		Dateadded:         published,
		IncludedWordLists: sourcePack.IncludedWordLists,
		StoreHashes:       StoreHashes{Full: hashes.Full, Source: hashes.Source, IncludedWordListItems: hashes.IncludedWordListItems},
	})

	file, _ := json.MarshalIndent(storeConfig, "", " ")

	_ = ioutil.WriteFile("store/store.json", file, 0644)

}

//func ConvertWordListsToWordListItems(sourcePack source.SourcePack) (returnItem []IncludedWordListItem) {
//
//	for _, i := range sourcePack.IncludedWordLists {
//
//		filepath := fmt.Sprintf("store/hash/%s/%s/%s", sourcePack.Type, sourcePack.Owner, sourcePack.Repo)
//		returnItem = append(returnItem, IncludedWordListItem{i, filepath})
//
//	}
//	return returnItem
//}

func Hash(objs ...interface{}) []byte {
	digester := crypto.SHA1.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return digester.Sum(nil)
}
