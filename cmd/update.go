package cmd

import (
	"fmt"
	"os"

	"github.com/micnncim/golaunch/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	// AppDir is application directory
	AppDir = "/Applications"
	// TmpDir is temporary directory
	TmpDir = "/tmp"
)

// UpdateCmd is command updating add data
func UpdateCmd() cli.Command {
	return cli.Command{
		Name:   "update",
		Usage:  "Update app data",
		Action: update,
	}
}

func update(c *cli.Context) error {
	dir := AppDir
	names, err := util.GetLocalAppNames(dir)
	if err != nil {
		return err
	}

	appDetails, _, err := util.FetchAppDetails(names)
	if err != nil {
		return err
	}

	dataDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDir += TmpDir
	os.RemoveAll(dataDir)
	os.Mkdir(dataDir, os.ModePerm)
	os.Chmod(dataDir, 0777)
	err = util.SaveAppDetails(appDetails, dataDir)
	if err != nil {
		return err
	}

	return nil
}
