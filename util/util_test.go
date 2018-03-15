package util

import (
	"fmt"
	"io/ioutil"
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

	officialApps, unofficialAppNames, err := GetAppItems(appNames)
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

	err = os.RemoveAll(appDir)
	if err != nil {
		log.Println(err)
	}
}

func TestSaveAppItem(t *testing.T) {
	appDir := "tmpAppDir"
	os.Mkdir(appDir, os.ModePerm)
	os.Chmod(appDir, 0777)
	evernote, err := os.Create(appDir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(appDir)
	if err != nil {
		t.Fatal(err)
	}

	apps, _, err := GetAppItems(appNames)
	if err != nil {
		t.Fatal(err)
	}

	appDataDir := "tmpAppdataDir"
	os.Mkdir(appDataDir, os.ModePerm)
	os.Chmod(appDataDir, 0777)
	err = SaveAppItem(apps[0], appDataDir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(appDir)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(appDataDir)
	if err != nil {
		log.Println(err)
	}
}

func TestReadAppData(t *testing.T) {
	appDir := "tmpAppDir"
	os.Mkdir(appDir, os.ModePerm)
	os.Chmod(appDir, 0777)
	evernote, err := os.Create(appDir + "/Evernote.app")
	if err != nil {
		log.Println(err)
	}
	defer evernote.Close()

	appNames, err := GetLocalAppNames(appDir)
	if err != nil {
		t.Fatal(err)
	}

	apps, _, err := GetAppItems(appNames)
	if err != nil {
		t.Fatal(err)
	}

	appDataDir := "tmpAppdataDir"
	os.Mkdir(appDataDir, os.ModePerm)
	os.Chmod(appDataDir, 0777)
	err = SaveAppItem(apps[0], appDataDir)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(appDataDir + "/Evernote")
	if err != nil {
		t.Fatal(err)
	}
	appItem, err := ReadAppData(b)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Read: %s", appItem)

	err = os.RemoveAll(appDir)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(appDataDir)
	if err != nil {
		log.Println(err)
	}
}
