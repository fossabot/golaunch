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
	Item struct {
		Name         string `json:"trackName"`
		Desc         string `json:"description"`
		ShortDesc    string
		Genres       []string `json:"genres"`
		PrimaryGenre string   `json:"primaryGenreName"`
	}

	Items []Item

	SearchResult struct {
		ResultCount int   `json:"resultCount"`
		Results     Items `json:"results"`
	}
)

const apiEndPoint = "https://itunes.apple.com/search"

// GetItems get app items by iTunes Search API
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
	ioutil.WriteFile(out, b, 0666)

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
		b, err := ioutil.ReadFile(dataDir + "/" + filename)
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

func (items Items) Render() ([]string, error) {
	max := 16

	for i := range items {
		items[i].Name = string([]rune(items[i].Name)[:max])
		items[i].ShortDesc = string([]rune(items[i].Desc)[:64])
	}

	var rows []string
	row := "{{.Name}}\t{{.PrimaryGenre}}\t{{.ShortDesc}}"

	for _, item := range items {
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, max+4, 0, '\t', 0)

		t, err := template.New("t").Parse(row)
		if err != nil {
			return nil, err
		}
		t.Execute(w, item)
		w.Flush()

		rows = append(rows, buf.String())
	}

	return rows, nil
}

func (item Item) SetShortDesc() {
	max := 64
	item.ShortDesc = string([]rune(item.Desc)[:max])
}

func (item Item) TrimDescNewLine() {
	item.Desc = strings.Trim(item.Desc, "\n")
}
