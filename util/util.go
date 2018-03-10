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
)

type (
	AppItem struct {
		Name             string   `json:"trackName"`
		description      string   `json:"description"`
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
			// log.Println(appName, "did't be found.")
			unofficialAppNames = append(unofficialAppNames, appName)
			continue
		}

		officialApps = append(officialApps, result.Results[0])

		time.Sleep(time.Second)
	}

	return officialApps, unofficialAppNames, nil
}
