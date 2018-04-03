package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"

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
	var apps AppDetails
	var unofficials []string
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Fetching app data..."
	s.Start()
	defer s.Stop()

	for _, appName := range appNames {
		v := url.Values{}
		v.Add("term", appName)
		v.Add("media", "software")
		v.Add("limit", "1")
		endPoint := fmt.Sprintf("%s?%s", apiEndPoint, v.Encode())
		res, err := http.Get(endPoint)
		if err != nil {
			log.Println(err)
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			continue
		}

		result := SearchResult{}
		if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		if len(result.Results) == 0 {
			unofficials = append(unofficials, appName)
			continue
		}

		app := result.Results[0]
		for _, appName := range appNames {
			if app.Name == appName {
				apps = append(apps, app)
			}
		}

		time.Sleep(time.Second)
	}

	return apps, unofficials, nil
}

// SaveAppDetail encodes appDetail => msgpack
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
	nameMax := 24
	descMax := 64

	for i := range appDetails {
		nameLen := len(string([]rune(appDetails[i].Name)))
		if nameLen > nameMax {
			appDetails[i].Name = string([]rune(appDetails[i].Name)[:nameMax])
		}
		appDetails[i].Name = strings.Trim(appDetails[i].Name, "\x00")

		descLen := len(string([]rune(appDetails[i].Desc)))
		if descLen > descMax {
			appDetails[i].Desc = string([]rune(appDetails[i].Desc)[:descMax])
		}
	}

	var rows []string
	tmpl := "{{.Name}}\t{{.PrimaryGenre}}\t{{.Desc}}"

	for _, appDetail := range appDetails {
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 32, 0, 4, ' ', 0)

		t, err := template.New("t").Parse(tmpl)
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
