package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type (
	SearchMap struct {
		TrackName        string   `json:"trackName"`
		Description      string   `json:"description"`
		Genres           []string `json:"genres"`
		PrimaryGenreName string   `json:"primaryGenreName"`
	}

	SearchResult struct {
		ResultCount int         `json:"resultCount"`
		Results     []SearchMap `json:"results"`
	}

	AppItem struct {
		Name         string
		Description  string
		Genres       []string
		PrimaryGenre string
	}

	AppItems []AppItem
)

const apiEndPoint = "https://itunes.apple.com/search"

func GetLocalAppNames() ([]string, error) {
	appDir := os.Getenv("GOLAUNCH_APP_DIR")
	fmt.Println(appDir)
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

func GetAppItems(appNames []string) (AppItems, error) {
	var appItems AppItems
	var count int

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
			log.Println(appName, "did't be found.")
			continue
		}

		if primaryName := strings.Fields(appName)[0]; !strings.Contains(result.Results[0].TrackName, primaryName) {
			log.Println(appName, "did't be found.")
			continue
		}

		count++
		fmt.Printf("%v\n", result.Results[0])

		time.Sleep(time.Second)
	}

	fmt.Printf("%d/%d found.\n", count, len(appNames))

	return appItems, nil
}
