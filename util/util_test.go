package util

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetLocalAppNames(t *testing.T) {
	appDir := "tmpAppDir"
	os.Mkdir(appDir, os.ModePerm)
	os.Chmod(appDir, 0777)
	sample1, err := os.Create(appDir + "/sample1.app")
	if err != nil {
		log.Println(err)
	}
	defer sample1.Close()
	sample2, err := os.Create(appDir + "/sample2.txt")
	if err != nil {
		log.Println(err)
	}
	defer sample2.Close()
	evernote, err := os.Create(appDir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(appDir)
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
	err = os.RemoveAll(appDir)
	if err != nil {
		log.Println(err)
	}
}

func TestGetAppItems(t *testing.T) {
	appDir := "tmpAppDir"
	os.Mkdir(appDir, os.ModePerm)
	os.Chmod(appDir, 0777)
	sample1, err := os.Create(appDir + "/sample1.app")
	if err != nil {
		log.Println(err)
	}
	defer sample1.Close()
	sample2, err := os.Create(appDir + "/sample2.txt")
	if err != nil {
		log.Println(err)
	}
	defer sample2.Close()
	evernote, err := os.Create(appDir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(appDir)
	if err != nil {
		t.Fatal(err)
	}

	appItems, err := GetAppItems(appNames)
	if err != nil {
		t.Fatal(err)
	}
	if len(appItems) != 1 {
		t.Fatal(err)
	}
	if appItems[0].Name != "Evernote" {
		t.Fatal(err)
	}

	err = os.RemoveAll(appDir)
	if err != nil {
		log.Println(err)
	}
}
