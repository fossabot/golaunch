package cmd

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"

	"github.com/micnncim/golaunch/util"
	"github.com/micnncim/interactive"
)

func Launch(c *cli.Context) error {
	dataDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDir += TmpDir

	appDetails, err := util.ReadAppDataFiles(dataDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	appDetailRows, err := appDetails.Render()
	if err != nil {
		fmt.Println(err)
		return err
	}

	command := "fzf"
	rows, err := interactive.Select(command, appDetailRows)
	if err != nil {
		fmt.Println(err)
		return err
	}

	app := strings.Split(rows[0], "\t")[0]
	if err := open.Run(app); err != nil {
		return err
	}

	return nil
}
