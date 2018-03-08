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
		t.Error("App names don't exist.")
	}

	fmt.Println(appNames)
}
