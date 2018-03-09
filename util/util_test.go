package util

import (
	"fmt"
	"testing"
)

func TestGetLocalAppNames(t *testing.T) {
	appNames, err := GetLocalAppNames()
	if err != nil {
		t.Fatal(err)
	}

	if len(appNames) == 0 {
		t.Fatal("App names don't exist.")
	}

	fmt.Println(appNames)
}

func TestGetAppItems(t *testing.T) {
	appNames, err := GetLocalAppNames()
	if err != nil {
		t.Fatal(err)
	}

	GetAppItems(appNames)
}
