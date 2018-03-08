package util

import (
	"io/ioutil"
	"regexp"
)

type AppData struct {
	Name         string
	Description  string
	Genres       []string
	PrimaryGenre string
}

const appDir = "/Applications"

func GetLocalAppNames() ([]string, error) {
	files, err := ioutil.ReadDir(appDir)
	if err != nil {
		return nil, err
	}

	var appNames []string
	for _, file := range files {
		fileName := file.Name()
		isApp, err := regexp.MatchString(`.*\.app`, fileName)
		if err != nil {
			return nil, err
		}
		if isApp {
			appNames = append(appNames, fileName)
		}
	}

	return appNames, nil
}
