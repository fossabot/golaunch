package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"

	"github.com/micnncim/golaunch/util"
	"github.com/micnncim/interactive"
)

func Launch(c *cli.Context) error {
	dir := "/Applications"
	names, err := util.GetLocalAppNames(dir)
	if err != nil {
		fmt.Printf("GetLocalAppNames: %s\n", err)
		return err
	}

	items, _, err := util.GetItems(names)
	if err != nil {
		fmt.Printf("GetItems: %s\n", err)
		return err
	}

	itemRows, err := items.Render()
	for _, itemRow := range itemRows {
		fmt.Println(itemRow)
	}
	if err != nil {
		fmt.Printf("Render: %s\n", err)
		return err
	}

	command := "fzf"
	rows, err := interactive.Select(command, itemRows)
	if err != nil {
		fmt.Printf("Select: %s\n", err)
		return err
	}

	for _, row := range rows {
		elements := strings.Split(row, "\t")
		app := elements[0] + ".app"
		cmd := exec.Command("open", "-a", app)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return nil
}
