package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

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

func Scan(prompt string) (string, error) {
	fmt.Print(prompt)

	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	return s.Text(), nil
}
