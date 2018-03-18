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
	dataDir := os.Getenv("HOME") + "/tmp"
	items, err := util.ReadAppDataFiles(dataDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	itemRows, err := items.Render()
	if err != nil {
		fmt.Println(err)
		return err
	}

	command := "fzf"
	rows, err := interactive.Select(command, itemRows)
	if err != nil {
		fmt.Println(err)
		return err
	}

	app := strings.Split(rows[0], "\t")[0]
	cmd := exec.Command("open", "-a", app)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
