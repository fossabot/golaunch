package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/vmihailenco/msgpack"
)

type (
	AppDetail struct {
		Name         string `json:"trackName"`
		Desc         string `json:"description"`
		ShortDesc    string
		Genres       []string `json:"genres"`
		PrimaryGenre string   `json:"primaryGenreName"`
	}

	AppDetails []AppDetail

	SearchResult struct {
		ResultCount int        `json:"resultCount"`
		Results     AppDetails `json:"results"`
	}
)

const apiEndPoint = "https://itunes.apple.com/search"

// FetchAppDetails get app appDetails by iTunes Search API
func FetchAppDetails(appNames []string) (AppDetails, []string, error) {
	var officialApps AppDetails
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

// SaveAppDetail encode appDetail => msgpack
func SaveAppDetail(appDetail AppDetail, dataDir string) error {
	b, err := msgpack.Marshal(&appDetail)
	if err != nil {
		return err
	}

	out := dataDir + "/" + appDetail.Name
	ioutil.WriteFile(out, b, 0666)

	return err
}

func SaveAppDetails(appDetails AppDetails, dataDir string) error {
	var err error
	for _, appDetail := range appDetails {
		err = SaveAppDetail(appDetail, dataDir)
		if err != nil {
			return err
		}
	}

	return err
}

func ReadAppData(data []byte) (AppDetail, error) {
	var appDetail AppDetail
	err := msgpack.Unmarshal(data, &appDetail)
	if err != nil {
		return appDetail, err
	}

	return appDetail, err
}

func ReadAppDataFiles(dataDir string) (AppDetails, error) {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	appDetails := AppDetails{}
	for _, file := range files {
		filename := file.Name()
		b, err := ioutil.ReadFile(dataDir + "/" + filename)
		if err != nil {
			log.Println(err)
			continue
		}
		appDetail, err := ReadAppData(b)
		if err != nil {
			log.Println(err)
			continue
		}
		appDetails = append(appDetails, appDetail)
	}

	return appDetails, nil
}

func (appDetails AppDetails) Render() ([]string, error) {
	max := 16

	for i := range appDetails {
		appDetails[i].Name = string([]rune(appDetails[i].Name)[:max])
		appDetails[i].ShortDesc = string([]rune(appDetails[i].Desc)[:64])
	}

	var rows []string
	row := "{{.Name}}\t{{.PrimaryGenre}}\t{{.ShortDesc}}"

	for _, appDetail := range appDetails {
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, max+4, 0, '\t', 0)

		t, err := template.New("t").Parse(row)
		if err != nil {
			return nil, err
		}
		t.Execute(w, appDetail)
		w.Flush()

		rows = append(rows, buf.String())
	}

	return rows, nil
}

func (appDetail AppDetail) SetShortDesc() {
	max := 64
	appDetail.ShortDesc = string([]rune(appDetail.Desc)[:max])
}

func (appDetail AppDetail) TrimDescNewLine() {
	appDetail.Desc = strings.Trim(appDetail.Desc, "\n")
}
