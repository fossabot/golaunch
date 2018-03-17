package util

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetLocalAppNames(t *testing.T) {
	dir := "tmpDir"
	os.Mkdir(dir, os.ModePerm)
	os.Chmod(dir, 0777)
	sample1, err := os.Create(dir + "/sample1.app")
	if err != nil {
		log.Println(err)
	}
	defer sample1.Close()
	sample2, err := os.Create(dir + "/sample2.txt")
	if err != nil {
		log.Println(err)
	}
	defer sample2.Close()
	evernote, err := os.Create(dir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(dir)
	// appNames, err := GetLocalAppNames("/Applications")
	if err != nil {
		t.Fatal(err)
	}

	if len(appNames) == 0 {
		t.Fatal("No app names exist.")
	}
	if len(appNames) != 2 {
		t.Fatal("The number of apps is wrong.")
	}

	fmt.Println(appNames)
	err = os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}
}
