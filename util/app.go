package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

type (
	Item struct {
		Name             string   `json:"trackName"`
		Description      string   `json:"description"`
		Genres           []string `json:"genres"`
		PrimaryGenreName string   `json:"primaryGenreName"`
	}

	Items []Item

	SearchResult struct {
		ResultCount int   `json:"resultCount"`
		Results     Items `json:"results"`
	}
)

const apiEndPoint = "https://itunes.apple.com/search"

func GetLocalAppNames(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var appNames []string
	for _, file := range files {
		fileName := file.Name()
		isApp, err := regexp.MatchString(`.*\.app`, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
		if isApp {
			fileName = strings.Trim(fileName, ".app")
			appNames = append(appNames, fileName)
		}
	}

	return appNames, nil
}

func GetItems(appNames []string) (Items, []string, error) {
	var officialApps Items
	var unofficialAppNames []string

	for _, appName := range appNames {
		v := url.Values{}
		v.Add("term", appName)
		v.Add("media", "software")
		v.Add("limit", "1")
		endPoint := apiEndPoint + "?" + v.Encode()
		fmt.Println(endPoint)
		res, err := http.Get(endPoint)
		if err != nil {
			log.Println(err)
			continue
		}
		defer res.Body.Close()

		result := SearchResult{}
		if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		if len(result.Results) == 0 {
			unofficialAppNames = append(unofficialAppNames, appName)
			continue
		}

		officialApps = append(officialApps, result.Results[0])

		time.Sleep(time.Second)
	}

	return officialApps, unofficialAppNames, nil
}

// SaveItem encode item => msgpack
func SaveItem(item Item, dataDir string) error {
	b, err := msgpack.Marshal(&item)
	if err != nil {
		return err
	}

	out := dataDir + "/" + item.Name
	ioutil.WriteFile(out, b, 0777)

	return err
}

func SaveItems(items Items, dataDir string) error {
	var err error
	for _, item := range items {
		err = SaveItem(item, dataDir)
		if err != nil {
			return err
		}
	}

	return err
}

func ReadAppData(data []byte) (Item, error) {
	var item Item
	err := msgpack.Unmarshal(data, &item)
	if err != nil {
		return item, err
	}

	return item, err
}

func ReadAppDataFiles(dataDir string) (Items, error) {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	items := Items{}
	for _, file := range files {
		filename := file.Name()
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		item, err := ReadAppData(b)
		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func Render() {

}
