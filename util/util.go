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
	AppItem struct {
		Name             string   `json:"trackName"`
		Description      string   `json:"description"`
		Genres           []string `json:"genres"`
		PrimaryGenreName string   `json:"primaryGenreName"`
	}

	AppItems []AppItem

	SearchResult struct {
		ResultCount int      `json:"resultCount"`
		Results     AppItems `json:"results"`
	}
)

const apiEndPoint = "https://itunes.apple.com/search"

func GetLocalAppNames(appDir string) ([]string, error) {
	files, err := ioutil.ReadDir(appDir)
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

func GetAppItems(appNames []string) (AppItems, []string, error) {
	var officialApps AppItems
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

// SaveAppItem encode appItem => msgpack
func SaveAppItem(appItem AppItem, appDataDir string) error {
	b, err := msgpack.Marshal(&appItem)
	if err != nil {
		return err
	}

	out := appDataDir + "/" + appItem.Name
	ioutil.WriteFile(out, b, 0777)

	return err
}

func SaveAppItems(appItems AppItems, appdataDir string) error {
	var err error
	for _, appItem := range appItems {
		err = SaveAppItem(appItem, appdataDir)
		if err != nil {
			return err
		}
	}

	return err
}

func ReadAppData(data []byte) (AppItem, error) {
	var appItem AppItem
	err := msgpack.Unmarshal(data, &appItem)
	if err != nil {
		return appItem, err
	}

	return appItem, err
}
