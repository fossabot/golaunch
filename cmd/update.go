package cmd

import (
	"fmt"
	"os"

	"github.com/micnncim/golaunch/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func UpdateCmd() cli.Command {
	return cli.Command{
		Name:   "update",
		Usage:  "Update app data",
		Action: update,
	}
}

func update(c *cli.Context) error {
	dir := "/Applications"
	names, err := util.GetLocalAppNames(dir)
	if err != nil {
		return err
	}

	items, _, err := util.GetItems(names)
	if err != nil {
		return err
	}

	dataDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDir += "/tmp"
	os.Mkdir(dataDir, os.ModePerm)
	os.Chmod(dataDir, 0777)
	err = util.SaveItems(items, dataDir)
	if err != nil {
		return err
	}

	return nil
}
