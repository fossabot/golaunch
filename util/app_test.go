package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGetItems(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}

	officialApps, unofficialAppNames, err := GetItems(appNames)
	fmt.Println("officialApps: ", officialApps)
	fmt.Println("unofficialAppNames:", unofficialAppNames)
	if err != nil {
		t.Fatal(err)
	}
	if len(officialApps) != 1 {
		t.Fatal("The number of apps is wrong.")
	}
	if officialApps[0].Name != "Evernote" {
		t.Fatal("The app name should be Evernote.")
	}

	err = os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}
}

func TestSaveItem(t *testing.T) {
	dir := "tmpDir"
	os.Mkdir(dir, os.ModePerm)
	os.Chmod(dir, 0777)
	evernote, err := os.Create(dir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(dir)
	if err != nil {
		t.Fatal(err)
	}

	apps, _, err := GetItems(appNames)
	if err != nil {
		t.Fatal(err)
	}

	dataDir := "tmpDataDir"
	os.Mkdir(dataDir, os.ModePerm)
	os.Chmod(dataDir, 0777)
	err = SaveItem(apps[0], dataDir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(dataDir)
	if err != nil {
		log.Println(err)
	}
}

func TestReadAppData(t *testing.T) {
	dir := "tmpDir"
	os.Mkdir(dir, os.ModePerm)
	os.Chmod(dir, 0777)
	evernote, err := os.Create(dir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(dir)
	if err != nil {
		t.Fatal(err)
	}

	apps, _, err := GetItems(appNames)
	if err != nil {
		t.Fatal(err)
	}

	dataDir := "tmpDataDir"
	os.Mkdir(dataDir, os.ModePerm)
	os.Chmod(dataDir, 0777)
	err = SaveItem(apps[0], dataDir)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(dataDir + "/Evernote")
	if err != nil {
		t.Fatal(err)
	}
	item, err := ReadAppData(b)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Read: %s", item)

	err = os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(dataDir)
	if err != nil {
		log.Println(err)
	}
}
