package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/micnncim/golaunch/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func AddCmd() cli.Command {
	return cli.Command{
		Name:   "add",
		Usage:  "Add app data",
		Action: add,
	}
}

func add(c *cli.Context) error {
	appName, err := util.Scan(color.CyanString("AppName> "))
	if err != nil {
		return err
	}

	desc, err := util.Scan(color.CyanString("Description> "))
	if err != nil {
		return err
	}

	genre, err := util.Scan(color.CyanString("Genre> "))
	if err != nil {
		return err
	}

	dataDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDir += DataDir

	err = util.SaveAppDetail(
		util.AppDetail{
			Name:         appName,
			Desc:         desc,
			ShortDesc:    desc,
			Genres:       []string{genre},
			PrimaryGenre: genre,
		},
		dataDir,
	)

	return err
}
